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
	"net/http"
	"sync"
	"time"

	"github.com/beefsack/go-rate"
)

const (
	retrieveErr     = "Error in fetching address"
	readErr         = "Error in reading the body of http responce; Error: "
	unmarshalErr    = "Error in unmarshaling the url details; Error: "
	notConnectedErr = "Given celebrities are not connected"
)

type Connection struct {
	person1          string
	person2          string
	bucketAddr       string
	connected        map[string]bool
	person2Mv        map[string]bool
	person2Detail    *Details
	urlToExplore     []url
	urlBeingExplored []url
	finish           chan []relation
	rw               sync.RWMutex
	wg               sync.WaitGroup
	rl               *rate.RateLimiter
}

func (c *Connection) Initialize(person1 string, person2 string, bucketAddr string) error {
	c.person1 = person1
	c.person2 = person2
	c.bucketAddr = bucketAddr
	c.connected = make(map[string]bool)
	c.person2Mv = make(map[string]bool)
	c.rl = rate.New(100, time.Second) // 200 times per second
	c.finish = make(chan []relation)
	return nil
}

func (c *Connection) foundMovie(url url, movie credit, name string) {
	var cred credit
	//fing this movie in person 2 detail
	for _, v := range c.person2Detail.Movies {
		if v.Url == movie.Url {
			cred = v
		}
	}
	rel := relation{movie.Name, name, movie.Role, c.person2, cred.Role}
	c.finish <- append(url.relation, rel)
}

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
func (c *Connection) findRelationShip() error {
	//swap urlToExplore and urlExplored
	temp := c.urlBeingExplored
	c.urlBeingExplored = c.urlToExplore
	c.urlToExplore = temp
	//if there is no url to be searched next
	//then that mean no connection possible
	fmt.Println("People to explore: ", len(c.urlBeingExplored))
	if len(c.urlBeingExplored) == 0 {
		return errors.New("Given celebrities are not connected")
	}
	for _, v := range c.urlBeingExplored {

		//for each entry get people they are connected to
		//get all the movie of this person
		c.wg.Add(1)
		go func(v url) {
			defer c.wg.Done()

			poi, err := c.fetchData(v.url)
			if err != nil {
				//return false, errors.New("error in retrieving address " + c.bucketAddr + c.person1 + "\n" + err.Error())
				return
			}
			//check wether movies of this person match that of person2
			for _, movie := range poi.Movies {
				if c.person2Mv[movie.Url] {
					c.foundMovie(v, movie, poi.Name)
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
						//return false, errors.New("error in retrieving address " + c.bucketAddr + c.person1 + "\n" + err.Error())
						return
					}

					for _, conn := range cnc.Cast {
						if c.isExplored(conn.Url) {
							continue
						}
						//new relation
						rel := relation{movie.Name, poi.Name, movie.Role, conn.Name, conn.Role}
						c.urlToExplore = append(c.urlToExplore, url{conn.Url, append(v.relation, rel)})
					}

					for _, conn := range cnc.Crew {
						if c.isExplored(conn.Url) {
							continue
						}
						//new connection
						rel := relation{movie.Name, poi.Name, movie.Role, conn.Name, conn.Role}

						c.urlToExplore = append(c.urlToExplore, url{conn.Url, append(v.relation, rel)})
					}
				}(movie)

			}
		}(v)
	}

	c.wg.Wait()
	return nil
}

func (c *Connection) GetRelationship() error {
	if c.person1 == c.person2 {
		//0 degree Separation
		c.finish <- nil
		return nil
	}

	//get details of both person
	p1Details, err := c.fetchData(c.person1)
	if err != nil {
		return errors.New("error in retrieving address " + c.bucketAddr + c.person1 + "\n" + err.Error())
	}

	p2Details, err := c.fetchData(c.person2)
	if err != nil {
		return errors.New("error in retrieving address " + c.bucketAddr + c.person2 + "\n" + err.Error())
	}

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

	c.urlToExplore = append(c.urlToExplore, url{c.person1, nil})
	c.connected[c.person1] = true

	for {
		err := c.findRelationShip()
		if err != nil {
			return err
		}
	}
}

//fetchData retrieve the data of a person or movie from the s3 bucket
func (c *Connection) fetchData(url string) (*Details, error) {
	//fetch the data
	c.rl.Wait()
	rs, err := http.Get(c.bucketAddr + url)
	if err != nil {
		fmt.Println("error in retrieving address " + c.bucketAddr + url + "\n" + err.Error())
		return nil, err
	}

	//read body of the data
	data, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		return nil, err
	}

	var detail Details
	err = json.Unmarshal(data, &detail)
	if err != nil {
		return nil, err
	}
	return &detail, nil
}
