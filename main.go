package main

import (
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var INVALID_VAL int = 100005

func main() {
	initData()
	findDegreesOfSeparation("amitabh-bachchan", "robert-de-niro")

	fmt.Println("Found answer ", iMoviesSeparationLevel)
}

func initData() {
	SetIsNodeVisited = make(map[string]bool)
}

type InfoNode struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Role string `json:"role"`
}

type MovieBuffResponseDataForPerson struct {
	URL    string     `json:"url"`
	Type   string     `json:"type"`
	Name   string     `json:"name"`
	Movies []InfoNode `json:"movies"`
}

type MovieBuffResponseDataForMovies struct {
	URL  string     `json:"url"`
	Type string     `json:"type"`
	Name string     `json:"name"`
	Cast []InfoNode `json:"cast"`
	Crew []InfoNode `json:"crew"`
}

// add hashset to find out nodes which we have already covered to escape loop
var SetIsNodeVisited map[string]bool

var iMoviesSeparationLevel int = 0

func findDegreesOfSeparation(sFirstPersonURL string, sSecondPersonURL string) (int, error) {
	sFunctionName := "findDegreesOfSeparation"

	if len(sFirstPersonURL) == 0 || len(sSecondPersonURL) == 0 {
		fmt.Println(sFunctionName, "Please provide valid inputs")
		return INVALID_VAL, errors.New("Invalid data provided")
	}

	ans := 0

	queue := list.New()

	addNodeToQueue(queue, InfoNode{URL: sFirstPersonURL})

	bIsCurrentPerson := true

	for false == SetIsNodeVisited[sSecondPersonURL] {

		nextQueue := list.New()

		for queue.Len() > 0 && false == SetIsNodeVisited[sSecondPersonURL] {
			frontElement := queue.Front()
			queue.Remove(frontElement)

			if nil == frontElement.Value {
				fmt.Println(sFunctionName, "Found invalid value during traversing at level", iMoviesSeparationLevel)
				return INVALID_VAL, errors.New(fmt.Sprintf("Found invalid value during traversing at level %d", iMoviesSeparationLevel))
			}

			node := frontElement.Value.(InfoNode)

			if len(node.URL) == 0 {
				fmt.Println(sFunctionName, "Found invalid value during traversing at level", iMoviesSeparationLevel)
				return INVALID_VAL, errors.New(fmt.Sprintf("Found invalid value during traversing at level %d", iMoviesSeparationLevel))
			}

			populateNeighbours(node.URL, nextQueue, bIsCurrentPerson)
		}

		// increase the movie count
		if false == bIsCurrentPerson {
			iMoviesSeparationLevel++
		}

		// current BFS done, use the new queue for next iteration
		queue = nextQueue

		// next time we will be iterating opposite
		bIsCurrentPerson = !bIsCurrentPerson
	}

	return ans, nil
}

func populateNeighbours(sMovieBuffURL string, queue *list.List, bIsCurrentPerson bool) error {
	sFunctionName := "populateNeighbours"

	if len(sMovieBuffURL) == 0 {
		fmt.Println(sFunctionName, "Invalid Moviebuff URL provided")
		return errors.New("Invalid Moviebuff URL provided")
	}

	sFormattedURL := fmt.Sprintf("http://data.moviebuff.com/%s", sMovieBuffURL)

	// Make GET request
	response, err := doHTTPRequest(sFormattedURL)

	if err != nil {
		fmt.Println(sFunctionName, "Error making the request:", err)
		return err
	}

	if bIsCurrentPerson {
		unmarshalAndPopulatePerson(response, queue)
	} else {
		unmarshalAndPopulateMovies(response, queue)
	}

	return nil
}

func unmarshalAndPopulatePerson(rawData []byte, queue *list.List) error {
	sFunctionName := "unmarshalAndPopulatePerson"

	// Unmarshal JSON into struct
	var responseData MovieBuffResponseDataForPerson
	err := json.Unmarshal(rawData, &responseData)
	if err != nil {
		fmt.Println(sFunctionName, "Error unmarshalling JSON:", err)
		return err
	}

	for _, node := range responseData.Movies {
		addNodeToQueue(queue, node)
	}

	return nil
}

func unmarshalAndPopulateMovies(rawData []byte, queue *list.List) error {
	sFunctionName := "unmarshalAndPopulatePerson"

	// Unmarshal JSON into struct
	var responseData MovieBuffResponseDataForMovies
	err := json.Unmarshal(rawData, &responseData)
	if err != nil {
		fmt.Println(sFunctionName, "Error unmarshalling JSON:", err)
		return err
	}

	for _, node := range responseData.Cast {
		addNodeToQueue(queue, node)
	}

	for _, node := range responseData.Crew {
		addNodeToQueue(queue, node)
	}

	return nil
}

func doHTTPRequest(sURL string) ([]byte, error) {
	sFunctionName := "doHTTPRequest"

	// Make GET request
	response, err := http.Get(sURL)
	if err != nil {
		fmt.Println(sFunctionName, "Error making the request:", err)
		return nil, err
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading the response body:", err)
		return nil, err
	}

	return body, nil
}

func addNodeToQueue(queue *list.List, node InfoNode) bool {

	if SetIsNodeVisited[node.URL] {
		return false
	}

	SetIsNodeVisited[node.URL] = true
	queue.PushBack(node)

	return true
}
