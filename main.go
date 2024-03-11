package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type movieMember interface {
	GetName() string
	GetURL() string
	GetRole() string
}

type movieMemberWrapper struct {
	Members []movieMember
}

func (m movieMemberWrapper) GetURL() string {
	if len(m.Members) > 0 {
		return m.Members[0].GetURL()
	}
	return ""
}

func (m movieMemberWrapper) AddMembers(members ...movieMember) {
	m.Members = append(m.Members, members...)
}

type dos struct {
	Movie          string
	Person1, Role1 string
	Person2, Role2 string
}

type personInfo struct {
	Name   string
	URL    string
	Movies []movie
	Cast   []cast
	Crew   []crew
}

type movie struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Type string `json:"type"`
	Role string `json:"role"`
}

type cast struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Role string `json:"role"`
}

type crew struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Role string `json:"role"`
}

func (c cast) GetName() string {
	return c.Name
}

func (c crew) GetName() string {
	return c.Name
}
func (c cast) GetURL() string {
	return c.URL
}

func (c crew) GetURL() string {
	return c.URL
}

func (c cast) GetRole() string {
	return c.Role
}
func (c crew) GetRole() string {
	return c.Role
}

const source = "http://data.moviebuff.com/"

type MovieBuff struct {
	Source, Destination string
	Person1, Person2    *personInfo
	Person2Movies       map[string]movie
	VisitedPerson       map[string]bool
	Visit               []string
	Visited             map[string]bool
	Link                map[string]dos
	Mutex               sync.Mutex
}

var (
	movieBuff    MovieBuff
)

func main() {
	args := os.Args[1:]

	if len(args) != 2 {
		log.Fatalln("Please provide two actor names")
	}

	if args[0] == args[1] {
		log.Fatalln("Degree of separation: 0 (Same actor)")
	}

	movieBuff.Person2Movies, movieBuff.Visited, movieBuff.Link, movieBuff.VisitedPerson = make(map[string]movie), make(map[string]bool), make(map[string]dos), make(map[string]bool)

	if err := processPersonData(args[0], args[1]); err != nil {
		log.Fatalln(err.Error())
	}

	t1 := time.Now()

	// Use a channel to signal when the processing is done
	done := make(chan struct{})
	defer close(done)

	// Print dots in a separate goroutine
	go printDots(done)

	// Process data and find the degree of separation
	dos, err := findDos()
	if err != nil {
		log.Fatalln(err.Error())
	}

	t2 := time.Now()

	// Print the results
	fmt.Printf("\nDegree of separation: %d\n\n", len(dos))
	for i, d := range dos {
		fmt.Printf("%d. Movie: %s\n   %s: %s\n   %s: %s\n\n", i+1, d.Movie, d.Role1, d.Person1, d.Role2, d.Person2)
	}

	fmt.Println("Total Time taken: ", t2.Sub(t1))
}

func processPersonData(person1, person2 string) error {
	pn1, err := fetchData(person1)
	if err != nil {
		return err
	}

	pn2, err := fetchData(person2)
	if err != nil {
		return err
	}

	movieBuff.Mutex.Lock()
	defer movieBuff.Mutex.Unlock()

	if len(pn1.Movies) > len(pn2.Movies) {
		movieBuff.Source, movieBuff.Destination = person2, person1
		movieBuff.Person1, movieBuff.Person2 = pn2, pn1
	} else {
		movieBuff.Source, movieBuff.Destination = person1, person2
		movieBuff.Person1, movieBuff.Person2 = pn1, pn2
	}

	for _, movie := range movieBuff.Person2.Movies {
		movieBuff.Person2Movies[movie.URL] = movie
	}

	movieBuff.Visit = append(movieBuff.Visit, movieBuff.Source)
	movieBuff.Visited[movieBuff.Source] = true

	return nil
}

func findDos() ([]dos, error) {
	var d []dos
	for {
		for _, person := range movieBuff.Visit {
			person1, err := fetchData(person)
			if err != nil {
				if isEOFError(err) {
					continue
				}
				continue
			}

			for _, p1movie := range person1.Movies {
				if movieBuff.Person2Movies[p1movie.URL].URL == p1movie.URL {
					if _, found := movieBuff.Link[person1.URL]; found {
						d = append(d, movieBuff.Link[person1.URL], dos{
							p1movie.Name, person1.Name, p1movie.Role,
							movieBuff.Person2.Name, movieBuff.Person2Movies[p1movie.URL].Role,
						})
					} else {
						d = append(d, dos{
							p1movie.Name, person1.Name, p1movie.Role,
							movieBuff.Person2.Name, movieBuff.Person2Movies[p1movie.URL].Role,
						})
					}
					return d, nil
				}
			}

			for _, p1movie := range person1.Movies {
				if movieBuff.Visited[p1movie.URL] {
					continue
				}

				movieBuff.Visited[p1movie.URL] = true

				p1moviedetail, err := fetchData(p1movie.URL)
				if err != nil {
					if isEOFError(err) {
						continue
					}
					continue
				}

				addToVisit(
					movieMemberWrapper{Members: castToMovieMember(movieBuff.Person1.Cast)},
					movieMemberWrapper{Members: crewToMovieMember(movieBuff.Person1.Crew)},
					movieMemberWrapper{Members: castToMovieMember(person1.Cast)},
					movieMemberWrapper{Members: crewToMovieMember(person1.Crew)},
				)

				for _, p1moviecast := range p1moviedetail.Cast {
					if movieBuff.Visited[p1moviecast.URL] {
						continue
					}

					movieBuff.Visited[p1moviecast.URL] = true
					movieBuff.Visit = append(movieBuff.Visit, p1moviecast.URL)
					movieBuff.Link[p1moviecast.URL] = dos{
						p1movie.Name, person1.Name, p1movie.Role, p1moviecast.Name, p1moviecast.Role,
					}
				}

				for _, p1moviecrew := range p1moviedetail.Crew {
					if movieBuff.Visited[p1moviecrew.URL] {
						continue
					}

					movieBuff.Visited[p1moviecrew.URL] = true
					movieBuff.Visit = append(movieBuff.Visit, p1moviecrew.URL)
					movieBuff.Link[p1moviecrew.URL] = dos{
						p1movie.Name, person1.Name, p1movie.Role, p1moviecrew.Name, p1moviecrew.Role,
					}
				}
			}
		}
	}
}

func fetchData(url string) (*personInfo, error) {
	resp, err := http.Get(source + url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var pi personInfo
	err = json.Unmarshal(result, &pi)
	if err != nil {
		return nil, errors.New("Please provide valid actor name input")
	}

	return &pi, nil
}

func addToVisit(members ...movieMemberWrapper) {
	for _, member := range members {
		for _, m := range member.Members {
			movieBuff.Visit = append(movieBuff.Visit, m.GetURL())
			movieBuff.Link[m.GetURL()] = dos{}
		}
	}
}

func isEOFError(err error) bool {
	return err != nil && err.Error() == "EOF"
}

func castToMovieMember(castList []cast) []movieMember {
	var members []movieMember
	for _, c := range castList {
		members = append(members, movieMember(c))
	}
	return members
}

func crewToMovieMember(crewList []crew) []movieMember {
	var members []movieMember
	for _, c := range crewList {
		members = append(members, movieMember(c))
	}
	return members
}

// Function to print dots until done signal is received
func printDots(done <-chan struct{}) {
	message := "Degree of separation is being calculated. Please wait"
	duration := 1 * time.Second // half second

	fmt.Print(message)

	for {
		select {
		case <-time.After(duration):
			fmt.Print(".")
		case <-done:
			return
		}
	}
}
