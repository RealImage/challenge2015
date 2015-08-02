/*
Purpose 	  : this file contain the functions
				that helps calculate the degree
				of connection between two celebrity
File Name	  : connection.go
Package		  : main
Date 		  : 01.08.2015
Author 		  : Mayank Patel
Date		Name		Modification
*/

//Degrees project main.go
//this project get the degree of connection between
//two celebrity and tells how they are connected
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/beefsack/go-rate"
)

//constants for errors
const (
	retrieveErr     = "Error in fetching address "
	readErr         = "Error in reading the body of http responce; Error: "
	unmarshalErr    = "Error in unmarshaling the url details; Error: "
	notConnectedErr = "Given celebrities are not connected"
	addrNotNilErr   = "Address cannot be nil"
	retrieveAddrErr = "error in retrieving address "
)

//Connection struct is used to find out the
//degree and relationship between two person
type Connection struct {
	person1          string            //person 1 url
	person2          string            //person 2 url
	config           *conf             //configuration
	connected        map[string]bool   //to store all already connected person and people
	person2Mv        map[string]bool   //to store all the movie os person 2
	person2Detail    *details          //person 2 detail
	urlToExplore     []person          //list of people to be explored in next iteration
	urlBeingExplored []person          //list of people being explored right now
	Finish           chan []relation   //to receive final result from go routines
	rw               sync.RWMutex      //mutax for connected map
	wg               sync.WaitGroup    //wait group to synchronize the go routine
	rl               *rate.RateLimiter //rate limiter
}

//Initialize initialized the connection struct.
//It takes person 1 and 2 url and configuration
func (c *Connection) Initialize(person1 string, person2 string, config *conf) error {
	c.person1 = person1
	c.person2 = person2

	if config.Address != "" {
		c.config = config
	} else {
		return errors.New(addrNotNilErr)
	}

	c.connected = make(map[string]bool)
	c.person2Mv = make(map[string]bool)

	if config.Limit > 0 {
		c.rl = rate.New(config.Limit, time.Second) // config.Limit times per second
	} else {
		log.Println("Invalid rate limit in the configuration file.")
		c.rl = rate.New(150, time.Second) //150 times per second
	}

	c.Finish = make(chan []relation)
	return nil
}

//isExplored checks if a given url is explored.
func (c *Connection) isExplored(url string) bool {
	c.rw.RLock()
	ok := c.connected[url]
	c.rw.RUnlock()
	if ok {
		return true
	}
	c.rw.Lock()
	c.connected[url] = true
	c.rw.Unlock()
	return false
}

//findRelationShip calculate the relationship between person1 and person2
func (c *Connection) findRelationShip() error {

	//swap urlToExplore and urlExplored
	temp := c.urlBeingExplored
	c.urlBeingExplored = c.urlToExplore
	c.urlToExplore = temp

	//if there is no url to be searched next
	//then that mean no connection possible
	fmt.Println("People to explore: ", len(c.urlBeingExplored))
	if len(c.urlBeingExplored) == 0 {
		return errors.New(notConnectedErr)
	}

	for _, persn := range c.urlBeingExplored {
		c.wg.Add(1)

		go func(p person) {
			defer c.wg.Done()

			//get all the details of this person of interest
			poi, err := c.fetchData(p.url)
			if err != nil {
				//log.Println(retrieveAddrErr + p.url + "\n" + err.Error())
				return
			}

			//check wether movies of this person match that of person2
			for _, movie := range poi.Movies {
				if c.person2Mv[movie.Url] {
					var cred credit
					//fing this movie in person 2 detail
					for _, v := range c.person2Detail.Movies {
						if v.Url == movie.Url {
							cred = v
						}
					}
					rel := relation{movie.Name, poi.Name, movie.Role, c.person2, cred.Role}
					//search complete finish the program
					c.Finish <- append(p.relation, rel)
					return
				}
			}

			//no movie matched. explore new people from these movies
			for _, movie := range poi.Movies {
				if c.isExplored(movie.Url) {
					continue
				}
				c.wg.Add(1)
				go func(movie credit) {
					defer c.wg.Done()
					//for each movies checkout the cast and crew
					cnc, err := c.fetchData(movie.Url)
					if err != nil {
						//return false, errors.New("error in retrieving address " + c.config.Address + c.person1 + "\n" + err.Error())
						return
					}

					for _, conn := range cnc.Cast {
						if c.isExplored(conn.Url) {
							continue
						}
						//new relation
						rel := relation{movie.Name, poi.Name, movie.Role, conn.Name, conn.Role}
						//append for next iteration
						c.urlToExplore = append(c.urlToExplore, person{conn.Url, append(p.relation, rel)})
					}

					for _, conn := range cnc.Crew {
						if c.isExplored(conn.Url) {
							continue
						}
						//new connection
						rel := relation{movie.Name, poi.Name, movie.Role, conn.Name, conn.Role}

						c.urlToExplore = append(c.urlToExplore, person{conn.Url, append(p.relation, rel)})
					}
				}(movie)
			}
		}(persn)
	}
	//wait for all go routine to finish
	c.wg.Wait()
	return nil
}

//GetConnection is the public function to get the
//degree of connection and relation between two movie star
func (c *Connection) GetConnection() ([]relation, error) {
	if c.person1 == c.person2 {
		//0 degree Separation
		return nil, nil
	}

	//get details of both person
	p1Details, err := c.fetchData(c.person1)
	if err != nil {
		return nil, errors.New(retrieveAddrErr + c.config.Address + c.person1 + "\n" + err.Error())
	}

	p2Details, err := c.fetchData(c.person2)
	if err != nil {
		return nil, errors.New(retrieveAddrErr + c.config.Address + c.person2 + "\n" + err.Error())
	}

	//start the search from person who have done less movie
	if len(p1Details.Movies) > len(p2Details.Movies) {
		temp := c.person1
		c.person1 = c.person2
		c.person2 = temp
		for _, v := range p1Details.Movies {
			c.person2Mv[v.Url] = true
		}
		c.person2Detail = p1Details
	} else {
		//save all the movie of person2. Save last(and most expensive) iteration
		for _, v := range p2Details.Movies {
			c.person2Mv[v.Url] = true
		}
		c.person2Detail = p2Details
	}
	c.urlToExplore = append(c.urlToExplore, person{c.person1, nil})
	c.connected[c.person1] = true
	go func() {
		for {
			err := c.findRelationShip()
			if err != nil {
				log.Fatalln(err.Error())
			}
		}
	}()
	return <-c.Finish, nil
}

//fetchData retrieve the data of a person or movie from the s3 bucket
func (c *Connection) fetchData(url string) (*details, error) {
	//fetch the data
	c.rl.Wait()
	rs, err := http.Get(c.config.Address + url)
	if err != nil {
		fmt.Println(retrieveErr + c.config.Address + url + "\n" + err.Error())
		return nil, err
	}

	//read body of the data
	data, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		return nil, err
	}

	var detail details
	err = json.Unmarshal(data, &detail)
	if err != nil {
		return nil, err
	}
	return &detail, nil
}
