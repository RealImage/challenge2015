package main

import (
	"encoding/json"
	"errors"
	//	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const BURST_DELAY = 15

type link struct {
	Source connection `json:"source"`
	Target connection `json:"target"`
	Movie  string     `json:"movie"`
}

type path struct {
	Links []link
}

func newPath(src string) path {
	links := make([]link, 0)
	initialPath := path{links}
	return initialPath.addFirstLink(src)
}

func (path path) String() string {
	pathToTarget := path.Links[1:]
	str := fmt.Sprintf("\nDegrees of separation : %d", len(pathToTarget))
	for i, link := range pathToTarget {
		str = fmt.Sprintf("%s\n\n%d. Movie: %s\n%s: %s \n%s: %s", str, i+1, link.Movie, link.Source.Role, link.Source.Name, link.Target.Role, link.Target.Name)
	}
	return str
}

func (path path) addFirstLink(src string) path {
	path.Links = append(path.Links, link{Target: connection{Url: src}})
	return path
}

func (path path) addLink(src *person, next connection, movie *movie) path {
	src_connection := connection{Url: src.Url, Name: src.Name, Role: movie.getRole(src.Url)}
	path.Links = append(path.Links, link{Source: src_connection, Target: next, Movie: movie.Name})
	return path
}

func (path path) lastPerson() string {
	return path.Links[len(path.Links)-1].Target.Url
}

type connection struct {
	Url  string `json:"url"`
	Name string `json:"name"`
	Role string `json:"role"`
}
type person struct {
	Url    string
	Name   string
	Movies []connection
}

type movie struct {
	Url  string
	Name string
	Cast []connection
	Crew []connection
}

type degreeResult struct {
	Degrees int    `json:"degrees"`
	Links   []link `json:"links"`
}

func (movie *movie) getPeopleInvolved() []connection {
	crew := movie.filterProductionCompany()
	return append(movie.Cast, crew...)
}

func (movie *movie) filterProductionCompany() []connection {
	filtered := make([]connection, 0)
	for _, item := range movie.Crew {
		if item.Role != "Production Company / Production" && item.Role != "Distributor / Distribution" {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func (movie *movie) getRole(personId string) string {
	for _, member := range movie.getPeopleInvolved() {
		if member.Url == personId {
			return member.Role
		}
	}
	fmt.Sprintf("Role could not be found for %s in movie %s", personId, movie.Url)
	return ""
}

var apiRootUrl = "http://data.moviebuff.com/"

type pathQueue struct {
	queue []path
}

func newPathQueue(initialPath path) *pathQueue {
	queue := &pathQueue{make([]path, 0)}
	queue.push(initialPath)
	return queue
}

func (pathQueue *pathQueue) push(path path) {
	pathQueue.queue = append(pathQueue.queue, path)
}

func (pathQueue *pathQueue) pop() path {
	first_item := pathQueue.queue[0]
	pathQueue.queue = pathQueue.queue[1:len(pathQueue.queue)]
	return first_item
}

func (pathQueue *pathQueue) isEmpty() bool {
	return len(pathQueue.queue) == 0
}

func main() {
	http.HandleFunc("/degree", func(w http.ResponseWriter, r *http.Request) {
		source := r.FormValue("source")
		target := r.FormValue("target")
		fmt.Printf("Request for %s and %s\n", source, target)
		path, _ := connect(source, target)
		bytes, err := json.Marshal(getJsonResponse(path))
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(path)
		w.Write(bytes)
	})
	http.ListenAndServe(":8080", nil)
}

func getJsonResponse(path path) degreeResult {
	return degreeResult{len(path.Links) - 1, path.Links[1:]}
}

func connect(source string, target string) (path, int) {
	path := newPath(source)
	fringe := newPathQueue(path)
	processedMovies := make(map[string]bool)
	processedPersons := make(map[string]bool)
	processedPersons[source] = true
	nodesExpanded := 0

	for {
		if fringe.isEmpty() {
			panic(errors.New("No solution"))
		}
		pathSoFar := fringe.pop()
		personSoFar := pathSoFar.lastPerson()
		//		fmt.Println("-->", personSoFar)
		if personSoFar == target {
			return pathSoFar, nodesExpanded
		}

		personChan := make(chan person)
		go fetchPerson(personSoFar, personChan)
		currentPersonInPath := <-personChan
		movies := fetchMovies(currentPersonInPath.Movies)
		for _, movie := range movies {
			if _, ok := processedMovies[movie.Url]; ok {
				continue
			}
			people := movie.getPeopleInvolved()
			for _, personNextLevel := range people {
				if personNextLevel.Url == target {
					return pathSoFar.addLink(&currentPersonInPath, personNextLevel, &movie), nodesExpanded
				}
				if _, ok := processedPersons[personNextLevel.Url]; !ok {
					nextPossiblePath := pathSoFar.addLink(&currentPersonInPath, personNextLevel, &movie)
					fringe.push(nextPossiblePath)
					processedPersons[personNextLevel.Url] = true
				}
				nodesExpanded++
			}
			processedMovies[movie.Url] = true
		}
	}
	return path, nodesExpanded
}

func fetchPerson(personId string, personChan chan person) {
	var body []byte
	var err error
	body, err = fetchResponse(apiRootUrl + personId)

	var person person
	err = json.Unmarshal(body, &person)
	if err != nil {
		fmt.Printf("parse error for %s, so ignoring\n", personId)
	}
	personChan <- person
}

func fetchPeople(peopleConnection []connection) []person {
	personChan := make(chan person)
	people := make([]person, 0)
	limiter := time.Tick(time.Millisecond * BURST_DELAY)
	for _, person := range peopleConnection {
		<-limiter
		go fetchPerson(person.Url, personChan)
	}
	for i := 0; i < len(peopleConnection); i++ {
		people = append(people, <-personChan)
	}
	return people
}

func fetchMovies(moviesConnection []connection) []movie {
	movieChan := make(chan movie)
	movies := make([]movie, 0)
	limiter := time.Tick(time.Millisecond * BURST_DELAY)
	for _, movie := range moviesConnection {
		<-limiter
		go fetchMovie(movie.Url, movieChan)
	}
	for i := 0; i < len(moviesConnection); i++ {
		movies = append(movies, <-movieChan)
	}
	return movies
}

func fetchMovie(movieId string, movieChannel chan movie) {
	var body []byte
	var err error
	body, err = fetchResponse(apiRootUrl + movieId)

	var movie movie
	err = json.Unmarshal(body, &movie)
	if err != nil {
		fmt.Printf("parse error for %s, so ignoring\n", movieId)
	}
	movieChannel <- movie
}

func fetchResponse(url string) ([]byte, error) {
	for {
		res, err := http.Get(url)
		if err != nil {
			time.Sleep(200)
			fmt.Println(err, ", so retrying")
			continue
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			perror(err)
		} else {
			return body, nil
		}
	}
}

func perror(err error) {
	if err != nil {
		panic(err)
	}
}
