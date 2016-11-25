/*
Purpose 	  : this file contain the functions that helps calculate the degree
				of connection between two celebrity
File Name	  : moviebuff.go
Package		  : moviebuff
Date 		  : 01.08.2015
Author 		  : Mayank Patel
Date		Name		Modification
*/

//moviebuff the degree of connection between two celebrity and tells how they
//are connected
package moviebuff

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/beefsack/go-rate"
)

//constants for errors
const (
	serverDownErr   = " Server is down or does not exist"
	notPersonErr    = "Not a person: "
	usrNotExistErr  = "does not exists: "
	retrieveErr     = "Error in fetching address "
	readErr         = "Error in reading the body of http responce; Error: "
	unmarshalErr    = "Error in unmarshaling the url details; Error: "
	notConnectedErr = "Given celebrities are not connected"
	addrNilErr      = "Address cannot be empty"
)

//Connection struct is used to find out the degree and relationship between two
//person
type Connection struct {
	person1, person2               string          //person 1 and person 2 url
	config                         *Conf           //configuration
	connected                      map[string]bool //to store all already connected person and movies
	p2Mv                           map[string]bool //to store all the movie os person 2
	p2Detail                       *details        //person 2 detail
	urlBeingExplored, urlToExplore []person        //list of people being explored right now and list of people to be explored in next iteration
	finish                         chan []Relation //to receive final result from go routines
	rw                             sync.RWMutex    //mutax for connected map
	wg                             sync.WaitGroup  //wait group to synchronize the go routine
	m                              sync.Mutex
	rl                             *rate.RateLimiter //rate limiter
	found                          bool
}

//Initialize initialized the connection struct. It takes person 1 and 2 url and
//configuration
func (c *Connection) Initialize(person1 string, person2 string, config *Conf) error {
	if config.Address == "" {
		return errors.New(addrNilErr)
	}

	//check if address is valid
	_, err := net.DialTimeout("tcp", strings.TrimLeft(strings.TrimRight(config.Address, "/"), `http://`)+":80", time.Second)
	if err != nil {
		return errors.New(config.Address + serverDownErr + err.Error())
	}
	c.config = config

	if config.RetryCount <= 0 {
		log.Println("Invalid connection-retry-count in the configuration file.")
		c.config.RetryCount = 10
	}

	if config.Limit > 0 {
		c.rl = rate.New(config.Limit, time.Second) // config.Limit times per second
	} else {
		log.Println("Invalid rate limit in the configuration file.")
		c.rl = rate.New(150, time.Second) //150 times per second
	}

	c.connected, c.p2Mv, c.finish = make(map[string]bool), make(map[string]bool), make(chan []Relation)
	c.person1, c.person2 = person1, person2

	return nil
}

//GetConnection is the public function to get the
//degree of connection and relation between two movie star
func (c *Connection) GetConnection() ([]Relation, error) {
	if c.person1 == c.person2 {
		//0 degree Separation
		return nil, nil
	}
	var p1Details, p2Details *details
	var err error

	//get details of both person
	if p1Details, err = c.fetchData(c.person1); err != nil {
		return nil, errors.New(usrNotExistErr + c.person1)
	}

	if p1Details.Typ != "Person" {
		return nil, errors.New(notPersonErr + c.person1)
	}

	if p2Details, err = c.fetchData(c.person2); err != nil {
		return nil, errors.New(c.person2 + usrNotExistErr)
	}

	if p1Details.Typ != "Person" {
		return nil, errors.New(notPersonErr + c.person2)
	}

	//start the search from person who have done less movie
	if len(p1Details.Movies) > len(p2Details.Movies) {
		c.person1, c.person2 = c.person2, c.person1
		for _, v := range p1Details.Movies {
			c.p2Mv[v.Url] = true
		}
		c.p2Detail = p1Details
	} else {
		//save all the movie of person2. Save last(and most expensive) iteration
		for _, v := range p2Details.Movies {
			c.p2Mv[v.Url] = true
		}
		c.p2Detail = p2Details
	}

	c.urlToExplore = append(c.urlToExplore, person{c.person1, nil})
	c.connected[c.person1] = true
	go func() {
		//keep looking for person 2 in bfs manner
		for {
			err := c.findRelationShip()
			if err != nil {
				log.Fatalln(err.Error())
			}
			if c.found {
				break
			}
		}
	}()
	return <-c.finish, nil
}

//findRelationShip calculate the relationship between person1 and person2
func (c *Connection) findRelationShip() error {
	//if there is no url to be searched next
	//then that mean no connection possible
	if len(c.urlToExplore) == 0 {
		return errors.New(notConnectedErr)
	}

	c.m.Lock()
	defer c.m.Unlock()

	//swap urlToExplore and urlExplored
	c.urlBeingExplored, c.urlToExplore = c.urlToExplore, c.urlBeingExplored

	//explore all the person to be explored in this depth
	c.wg.Add(len(c.urlBeingExplored))
	for _, persn := range c.urlBeingExplored {
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
				if c.p2Mv[movie.Url] {
					var cred credit

					//fing this movie in person 2 detail
					for _, v := range c.p2Detail.Movies {
						if v.Url == movie.Url {
							cred = v
						}
					}
					rel := Relation{movie.Name, poi.Name, movie.Role, c.p2Detail.Name, cred.Role}

					//search complete. finish the program
					c.finish <- append(p.relation, rel)
					c.found = true
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
					if err != nil { //return false, errors.New("error in retrieving address " + c.config.Address + c.person1 + "\n" + err.Error())
						return
					}

					for _, conn := range cnc.Cast {
						if c.isExplored(conn.Url) {
							continue
						}

						//append for next iteration
						c.urlToExplore = append(c.urlToExplore, person{conn.Url, append(p.relation, Relation{movie.Name, poi.Name, movie.Role, conn.Name, conn.Role})})
					}

					for _, conn := range cnc.Crew {
						if c.isExplored(conn.Url) {
							continue
						}
						//append for next iteration
						c.urlToExplore = append(c.urlToExplore, person{conn.Url, append(p.relation, Relation{movie.Name, poi.Name, movie.Role, conn.Name, conn.Role})})
					}
				}(movie)
			}
		}(persn)
	}
	//wait for all go routine to finish
	c.wg.Wait()
	return nil
}

//fetchData retrieve the data of a person or movie from the s3 bucket
func (c *Connection) fetchData(url string) (details *details, err error) {
	//fetch the data
	c.rl.Wait()
	var rs *http.Response
	if rs, err = http.Get(c.config.Address + url); err != nil {
		for i := 0; i < c.config.RetryCount; i++ {
			//fmt.Println("trying again Error: ", i, err.Error())
			c.rl.Wait()
			if rs, err = http.Get(c.config.Address + url); err == nil {
				break
			}
			if strings.Contains(err.Error(), "too many open files") {
				//throttle the access to fetch data to cool down a bit
				for j := 0; j < c.config.Limit/4; j++ {
					c.rl.Wait()
				}
			}
		}
		if err != nil {
			log.Println(retrieveErr + c.config.Address + url + "\n" + err.Error())
			return nil, err
		}
	}
	defer rs.Body.Close()

	//read body of the data
	var data []byte
	if data, err = ioutil.ReadAll(rs.Body); err != nil {
		return nil, err
	}

	//unmarshal data into detail
	if err = json.Unmarshal(data, &details); err != nil {
		return nil, err
	}
	return details, nil
}

//isExplored checks if a given url is already explored.
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
