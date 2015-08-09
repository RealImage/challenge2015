//Copyright 2015 Mahendra Kathirvel. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type cast struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	Role string `json:"role"`
}

type crew struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	Role string `json:"role"`
}

type movie struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	Type string `json:"type"`
	Role string `json:"role"`
}

type dos struct {
	movie          string
	person1, role1 string
	person2, role2 string
}

// Used to parse the moviebuff json response
type info struct {
	Name   string  `json:"name"`
	Url    string  `json:"url"`
	Type   string  `json:"type"`
	Movies []movie `json:"movies"`
	Cast   []cast  `json:"cast"`
	Crew   []crew  `json:"crew"`
}

// Primary struct to hold the execution data
type Moviebuff struct {
	source        string
	destination   string
	person1       *info
	person2       *info
	p2Movies      map[string]movie
	visitedPerson map[string]bool
	visit         []string
	visited       map[string]bool
	link          map[string]dos
}

// Moviebuff data url
const source = "http://data.moviebuff.com/"

// Global variables
var (
	buff         Moviebuff
	totalRequest uint
)

func main() {

	// Parsing arguments
	args := os.Args[1:]

	if len(args) != 2 {
		log.Fatalln("\nPlease check arguments: Need two of them")
	}

	if args[0] == args[1] {
		log.Fatalln("\nDegree of seperation: 0")
	}

	// Init variable
	buff.p2Movies, buff.visited, buff.link, buff.visitedPerson = make(map[string]movie), make(map[string]bool), make(map[string]dos), make(map[string]bool)

	// Processing person data to start the execution
	if err := processPersonData(args[0], args[1]); err != nil {
		log.Fatalln(err.Error())
	}

	t1 := time.Now()

	// Find the relation between given person
	dos, err := findDos()
	if err != nil {
		log.Fatalln(err.Error())
	}

	t2 := time.Now()

	// Print the result
	fmt.Printf("\nDegree of separation: %d\n\n", len(dos))
	for i, d := range dos {
		fmt.Printf("%d. Movie: %s\n   %s: %s\n   %s: %s\n\n", i+1, d.movie, d.role1, d.person1, d.role2, d.person2)
	}

	// Optional stats
	fmt.Println("Total HTTP request sent: ", totalRequest)
	fmt.Println("Time taken: ", t2.Sub(t1))
}

// Fetch and store the data in a global variable
func processPersonData(person1, person2 string) error {

	pn1, err := fetchData(person1)
	if err != nil {
		return err
	}

	pn2, err := fetchData(person2)
	if err != nil {
		return err
	}

	if len(pn1.Movies) > len(pn2.Movies) {
		buff.source, buff.destination = person2, person1
		buff.person1, buff.person2 = pn2, pn1
	} else {
		buff.source, buff.destination = person1, person2
		buff.person1, buff.person2 = pn1, pn2
	}

	for _, movie := range buff.person2.Movies {
		buff.p2Movies[movie.Url] = movie
	}

	buff.visit = append(buff.visit, buff.source)
	buff.visited[buff.source] = true

	return nil
}

// Apply BFS to expore the each node
func findDos() ([]dos, error) {

	var d []dos
	for true {

		//fmt.Printf("Visited Person: %v, %d\n\n", buff.visitedPerson, len(buff.visitedPerson))

		for _, person := range buff.visit {

			/*fmt.Printf("%s\n\n", person)
			if buff.visitedPerson[person] {
				continue
			}
			buff.visitedPerson[person] = true
			*/

			person1, err := fetchData(person)
			if err != nil {
				if strings.Contains(err.Error(), "looking for beginning of value") {
					continue
				}
				return nil, err
			}

			for _, p1movie := range person1.Movies {
				if buff.p2Movies[p1movie.Url].Url == p1movie.Url {
					if _, found := buff.link[person1.Url]; found {
						d = append(d, buff.link[person1.Url], dos{p1movie.Name, person1.Name, p1movie.Role, buff.person2.Name, buff.p2Movies[p1movie.Url].Role})
					} else {
						d = append(d, dos{p1movie.Name, person1.Name, p1movie.Role, buff.person2.Name, buff.p2Movies[p1movie.Url].Role})
					}
					return d, nil
				}
			}

			// Find new nodes to continue searching
			for _, p1movie := range person1.Movies {

				if buff.visited[p1movie.Url] {
					continue
				}

				buff.visited[p1movie.Url] = true

				p1moviedetail, err := fetchData(p1movie.Url)
				if err != nil {
					if strings.Contains(err.Error(), "looking for beginning of value") {
						continue
					}
					return nil, err
				}

				for _, p1moviecast := range p1moviedetail.Cast {

					if buff.visited[p1moviecast.Url] {
						continue
					}

					buff.visited[p1moviecast.Url] = true
					buff.visit = append(buff.visit, p1moviecast.Url)
					buff.link[p1moviecast.Url] = dos{p1movie.Name, person1.Name, p1movie.Role, p1moviecast.Name, p1moviecast.Role}
				}

				for _, p1moviecrew := range p1moviedetail.Crew {

					if buff.visited[p1moviecrew.Url] {
						continue
					}

					buff.visited[p1moviecrew.Url] = true
					buff.visit = append(buff.visit, p1moviecrew.Url)
					buff.link[p1moviecrew.Url] = dos{p1movie.Name, person1.Name, p1movie.Role, p1moviecrew.Name, p1moviecrew.Role}
				}

			}
		}

		//fmt.Printf("Visit: %v, %d\n\n", buff.visit, len(buff.visit))
		//fmt.Printf("Visited: %v, %d\n\n", buff.visited, len(buff.visited))

	}

	return d, nil
}

// Fetch and parse the incoming json response
func fetchData(url string) (*info, error) {

	// Throttle the data request rate
	time.Sleep(100 * time.Millisecond)

	resp, _ := http.Get(source + url)
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var i info
	err = json.Unmarshal(result, &i)
	if err != nil {
		//log.Println(err, err.Error())
		return nil, err
	}

	totalRequest++
	//fmt.Printf("%v\n\n", i)
	return &i, nil
}
