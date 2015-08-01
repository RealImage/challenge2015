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
	urlToExplore     []url
	urlBeingExplored []url
	result           []relation
	count            int
	rw               sync.RWMutex
	wg               sync.WaitGroup
}

func (c *Connection) Initialize(person1 string, person2 string, bucketAddr string) error {
	c.person1 = person1
	c.person2 = person2
	c.bucketAddr = bucketAddr
	c.connected = make(map[string]bool)
	return nil
}

func (c *Connection) findRelationShip() error {
	//swap urlToExplore and urlExplored
	temp := c.urlBeingExplored
	c.urlBeingExplored = c.urlToExplore
	c.urlToExplore = temp
	//if there is no url to be searched next
	//then that mean no connection possible
	fmt.Println("dbbndfbndflbkndfkbndfb", len(c.urlBeingExplored))
	if len(c.urlBeingExplored) == 0 {
		return errors.New("Given celebrities are not connected")
	}

	for _, v := range c.urlBeingExplored {
		c.wg.Add(1)
		//for each entry get people they are connected to
		//get all the movie of this person
		go func(v url) {
			defer c.wg.Done()
			poi, err := c.fetchData(v.url)
			if err != nil {
				//return false, errors.New("error in retrieving address " + c.bucketAddr + c.person1 + "\n" + err.Error())
				return
			}

			for _, movie := range poi.Movies {
				c.rw.RLock()
				ok := c.connected[movie.Url]
				c.rw.RUnlock()

				if ok {
					continue
				}
				c.rw.Lock()
				c.connected[movie.Url] = true
				c.rw.Unlock()
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
						c.rw.RLock()
						ok := c.connected[conn.Url]
						c.rw.RUnlock()
						if ok {
							continue
						}
						//new connection
						var rel relation
						rel.movie = movie.Name
						rel.person1 = poi.Name
						rel.role1 = movie.Role
						rel.person2 = conn.Name
						rel.role2 = conn.Role
						if conn.Url == c.person2 {
							c.result = append(v.relation, rel)
							return
						}
						c.rw.Lock()
						c.connected[conn.Url] = true
						c.rw.Unlock()
						c.urlToExplore = append(c.urlToExplore, url{conn.Url, append(v.relation, rel)})
					}

					for _, conn := range cnc.Crew {
						c.rw.RLock()
						ok := c.connected[conn.Url]
						c.rw.RUnlock()
						if ok {
							continue
						}
						//new connection
						var rel relation
						rel.movie = movie.Name
						rel.person1 = poi.Name
						rel.role1 = movie.Role
						rel.person2 = conn.Name
						rel.role2 = conn.Role
						if conn.Url == c.person2 {
							c.result = append(v.relation, rel)
							return
						}
						c.rw.Lock()
						c.connected[conn.Url] = true
						c.rw.Unlock()
						c.urlToExplore = append(c.urlToExplore, url{conn.Url, append(v.relation, rel)})
					}
				}(movie)

			}
		}(v)
	}
	c.wg.Wait()
	return nil
}

func (c *Connection) GetRelationship() ([]relation, error) {
	if c.person1 == c.person2 {
		//0 degree Separation
		return nil, nil
	}

	//get details of both person
	p1Details, err := c.fetchData(c.person1)
	if err != nil {
		return nil, errors.New("error in retrieving address " + c.bucketAddr + c.person1 + "\n" + err.Error())
	}

	p2Details, err := c.fetchData(c.person2)
	if err != nil {
		return nil, errors.New("error in retrieving address " + c.bucketAddr + c.person2 + "\n" + err.Error())
	}

	if len(p1Details.Movies) > len(p2Details.Movies) {
		temp := c.person1
		c.person1 = c.person2
		c.person2 = temp
	}

	c.urlToExplore = append(c.urlToExplore, url{c.person1, nil})
	c.connected[c.person1] = true

	for {
		err := c.findRelationShip()
		if err != nil || c.result != nil {
			return c.result, err
		}
	}
}

//fetchData retrieve the data of a person or movie from the s3 bucket
func (c *Connection) fetchData(url string) (*Details, error) {
	t1 := time.Now()
	//fetch the data
	c.count++
	rs, err := http.Get(c.bucketAddr + url)
	if err != nil {
		fmt.Println("error in retrieving address " + c.bucketAddr + url + "\n" + err.Error())
		return nil, err
	}
	fmt.Println("Fetch :: ", time.Since(t1), url)

	//t2 := time.Now()
	//read body of the data
	data, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		return nil, err
	}
	//fmt.Println("read :: ", time.Since(t2))
	//t3 := time.Now()
	var detail Details
	err = json.Unmarshal(data, &detail)
	if err != nil {
		return nil, err
	}
	//fmt.Println("unmarshel :: ", time.Since(t3))

	return &detail, nil
}
