package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"
)

type Person struct {
	MovieBuffData
	Movies []MovieBuffDataArray `json:"movies"`
}

type Movie struct {
	MovieBuffData
	Cast []MovieBuffDataArray `json:"cast"`
	Crew []MovieBuffDataArray `json:"crew"`
}

type MovieBuffData struct {
	URL  string `json:"url"`
	Type string `json:"type"`
	Name string `json:"name"`
}

type MovieBuffDataArray struct {
	URL  string `json:"url"`
	Role string `json:"type"`
	Name string `json:"name"`
}

type TreeNode struct {
	URL    string
	Degree int
	Path   []string
}

const (
	PERSON_TYPE = iota
	MOVIE_TYPE
	MOVIE_BUFF_SERVER_URL = "https://data.moviebuff.com/"
)

type urlDataResponse interface {
	typeOfURL() string
}

func (person Person) typeOfURL() string {
	return person.Type
}

func (movie Movie) typeOfURL() string {
	return movie.Type
}

// This function retrieves data from the movieBuff server, with the specific data source being determined by the type of URL provided.
func fetchURLData(url string, urltype int) (urlDataResponse, error) {
	resp, err := http.Get(MOVIE_BUFF_SERVER_URL + url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if urltype == PERSON_TYPE {
		actor := Person{}
		if err := json.Unmarshal(responseData, &actor); err != nil {
			return nil, err
		}
		return actor, nil
	}

	movie := Movie{}
	if err := json.Unmarshal(responseData, &movie); err != nil {
		return nil, err
	}
	return movie, nil

}

func findDegreesOfSeparation(actor1, actor2 string) (*TreeNode, error) {
	startQueue := TreeNode{
		URL:    actor1,
		Degree: 0,
		Path:   []string{actor1},
	}

	//create queue to stored person and movie data
	queue := make([]*TreeNode, 0, 1)
	queue = append(queue, &startQueue)

	urlType := PERSON_TYPE
	visitedNodes := make(map[string]bool)

	//create mutex to lock and unlock of objects.
	var mu sync.Mutex

	//Create waitgroup to wait for a goroutine
	var wg sync.WaitGroup
	numofCore := runtime.NumCPU() //Get the number of core of cpu where this program is running
	corePool := make(chan struct{}, numofCore)

	//To process the queue until the queue is not empty or until we found second actor.
	for len(queue) > 0 {
		levelSize := len(queue)

		for i := 0; i < levelSize; i++ {

			URL := queue[i].URL
			degree := queue[i].Degree
			path := queue[i].Path

			if PERSON_TYPE == urlType {
				degree = degree + 1
			}

			if URL == actor2 {
				queue[i].Degree = degree
				return queue[i], nil
			}

			corePool <- struct{}{} // acquire a cpu cores to control the number of concurrent requests
			wg.Add(1)

			go func(url string) {
				defer func() {
					<-corePool
					wg.Done()
				}()
				if response, err := fetchURLData(url, urlType); err == nil {
					mu.Lock()
					visitedNodes[URL] = true
					mu.Unlock()
					urlType := response.typeOfURL()
					if urlType == "Person" { // Procced this if block when the urlType is Person.
						resp := response.(Person)
						for _, person := range resp.Movies {
							fmt.Printf("\nThis is person: %s", person.Name)
							path := make([]string, len(path))
							copy(path, path)
							mu.Lock()
							if !visitedNodes[person.URL] {
								URL = person.URL
								path = append(path, person.URL)
								q := TreeNode{
									URL:    URL,
									Degree: degree,
									Path:   path,
								}
								queue = append(queue, &q)
							}
							mu.Unlock()
						}
					} else if urlType == "Movie" { // Procced this if block when the urlType is Movie.
						resp := response.(Movie)
						for _, movie := range resp.Cast { //Proccess the Cast data
							fmt.Printf("\nThis is movie: %s", movie.Name)
							path := make([]string, len(path))
							copy(path, path)
							mu.Lock()
							if !visitedNodes[movie.URL] {
								URL = movie.URL
								path = append(path, movie.URL)
								q := TreeNode{
									URL:    URL,
									Degree: degree,
									Path:   path,
								}
								queue = append(queue, &q)
							}
							mu.Unlock()
						}

					}
				}
			}(URL)
		}

		wg.Wait()
		urlType = 1 - urlType
		queue = queue[levelSize:]
	}

	// Actors are not connected
	return nil, fmt.Errorf("No path found")
}

func main() {
	// Sample data for actors Amitabh Bachchan and Robert De Niro

	startTime := time.Now()
	if len(os.Args) < 3 {
		fmt.Println("Error: Invalid number of arguments. Please provide two actor names.")
		fmt.Println("Usage: degrees <actor-name> <actor-name>")
		return
	}

	actor1 := os.Args[1]
	actor2 := os.Args[2]

	result, err := findDegreesOfSeparation(actor1, actor2)
	if err != nil {
		fmt.Errorf("Error from findDegreesOfSeparation: %+v", err)
		return
	}

	endTime := time.Now()

	timeDifference := endTime.Sub(startTime)
	fmt.Printf("Time Difference: %s\n", timeDifference)
	if result == nil {
		fmt.Println("Actors are not connected.")
	} else {
		fmt.Println("Degree", result.Degree)
	}
}
