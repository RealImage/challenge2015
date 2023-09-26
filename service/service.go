package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	// "os"
	"runtime"
	"sync"
)

type queue struct {
	URL    string
	degree int
	path   []string
}

type BufferData struct {
	URL  string `json:"url"`
	Role string `json:"role"`
	Name string `json:"name"`
}

type BufferMetaData struct {
	URL  string `json:"url"`
	Type string `json:"type"`
	Name string `json:"name"`
}

type Person struct {
	BufferMetaData
	Movies []BufferData `json:"movies"`
}

type Movie struct {
	BufferMetaData
	Cast []BufferData `json:"cast"`
}

type FetchResponse interface {
	urlType() string
}

func (person Person) urlType() string {
    return person.Type
}

func (movie Movie) urlType() string {
    return movie.Type
}

const (
	CategoryPerson = iota
	CategoryMovie
)

func fetchAPIData(url string, URLCategiry int) (FetchResponse, error) {
	// fmt.Println("Fetching data for URL: ", "https://data.moviebuff.com/" + url)
	resp, err := http.Get("https://data.moviebuff.com/" + url)
	if err != nil {
		return nil, err
	}	
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	// fmt.Println("Response data: ", string(responseData))
	if err != nil {
		return nil, err
	}

	if URLCategiry == CategoryPerson {
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

func FindDegreesOfSeparation(firstActorName string, secondActorName string) (*queue, error) {
	
	// used to keep track of URLs we've already visited, so we don't visit them again
	visitedURLs := make(map[string]bool)
	startQueue := queue{URL: firstActorName, degree: 0, path: []string{firstActorName}}
	BFSQueue := make([]*queue, 0, 1)
	BFSQueue = append(BFSQueue, &startQueue)
	var mu sync.Mutex
	var wg sync.WaitGroup

	URLCategory := CategoryPerson
	
	totalWorkers := runtime.NumCPU()
	fmt.Println("Total workers: ", totalWorkers)
	workerPool := make(chan struct{}, totalWorkers)

	for len(BFSQueue) > 0 {
		length := len(BFSQueue)
		for i := 0; i < length; i++ {
			URL := BFSQueue[i].URL
			degree := BFSQueue[i].degree
			path := BFSQueue[i].path

			if URLCategory == CategoryPerson {degree = degree + 1}

			if URL == secondActorName {
				// fmt.Println("Found the path!")
				BFSQueue[i].degree = degree
				return BFSQueue[i], nil
			}

			workerPool <- struct{}{} // acquire a worker pool to control the number of concurrent workers
			wg.Add(1)

			go func(url string) {
				defer func() {
					<-workerPool // release the worker pool
					wg.Done()
				}()
				if response, err := fetchAPIData(url, URLCategory); err == nil {
					// fmt.Println("Response: ", response)
					// os.Exit(1)
					mu.Lock()
					visitedURLs[URL] = true
					mu.Unlock()
					urlCategory := response.urlType()
					// fmt.Println("URL category: ", urlCategory)
					if urlCategory == "Person" {
						resp := response.(Person)
						for _, person := range resp.Movies {
							path := make([]string, len(path))
							copy(path, path)
							mu.Lock()
							if !visitedURLs[person.URL] {
								URL = person.URL
								path = append(path, person.URL)
								q := queue{
									URL:    URL,
									degree: degree,
									path:   path,
								}
								BFSQueue = append(BFSQueue, &q)
							}
							mu.Unlock()
						}
					} else if urlCategory == "Movie" {
						resp := response.(Movie)
						for _, movie := range resp.Cast {
							path := make([]string, len(path))
							copy(path, path)
							mu.Lock()
							if !visitedURLs[movie.URL] {
								URL = movie.URL
								path = append(path, movie.URL)
								q := queue{
									URL:    URL,
									degree: degree,
									path:   path,
								}
								BFSQueue = append(BFSQueue, &q)
							}
							mu.Unlock()
						} 
					} 
				}
			} (URL)
		}
		wg.Wait()
		URLCategory = 1 - URLCategory
		BFSQueue = BFSQueue[length:]
	}
	return nil, fmt.Errorf("No path found")
}

func DisplayResult(firstActorName string, secondActorName string)  {
	// fmt.Println("Displaying result")
	result, err := FindDegreesOfSeparation(firstActorName, secondActorName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Degrees of separation: ", result.degree)
	fmt.Println("Path: ", strings.Join(result.path, " -> "))
}