// TODO - more refined error handling to check whether the provided arguments are actually person URL
// TODO - format the result.

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
)

const (
	TYPE_PERSON = iota
	TYPE_MOVIE
)

type MovieBuffData struct {
	URL  string `json:"url"`
	Role string `json:"role"`
	Name string `json:"name"`
}

type MovieBuffMetaData struct {
	URL  string `json:"url"`
	Type string `json:"type"`
	Name string `json:"name"`
}

type PersonData struct {
	MovieBuffMetaData
	Movies []MovieBuffData `json:"movies"`
}

type MovieData struct {
	MovieBuffMetaData
	Cast []MovieBuffData `json:"cast"`
}

type queueData struct {
	pathToCurrentURL []string
	currentURL       string
	degree           int
}

type MovieBuffResponse interface {
	urlType() string
}

func (person PersonData) urlType() string {
	return person.Type
}

func (movie MovieData) urlType() string {
	return movie.Type
}

func newQueue(url string) queueData {
	path := make([]string, 0, 1)
	path = append(path, url)

	return queueData{
		currentURL:       url,
		degree:           0,
		pathToCurrentURL: path,
	}
}

func fetchMoviebuffData(url string, urlType int) (MovieBuffResponse, error) {
	resp, err := http.Get("https://data.moviebuff.com/" + url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if urlType == TYPE_PERSON {
		person := PersonData{}
		if err := json.Unmarshal(data, &person); err != nil {
			return nil, err
		}
		return person, nil
	}

	movie := MovieData{}
	if err := json.Unmarshal(data, &movie); err != nil {
		return nil, err
	}
	return movie, nil

}

// we use breath first algorithm to find the degrees of separation.
func findDegreesOfSeparation(startURL, endURL string) (*queueData, error) {
	visited := make(map[string]bool) // to help prevent visiting the same URL again
	qData := newQueue(startURL)
	queue := make([]queueData, 0, 1)
	queue = append(queue, qData)
	urlType := TYPE_PERSON
	var mu sync.Mutex // Mutex to protect visited map
	var wg sync.WaitGroup

	// Setting number of concurrent workers as per the number of available logical cores
	numWorkers := runtime.NumCPU() //this number can be more than the number of logical cores. Need trial & error to fix the correct number

	semaphore := make(chan struct{}, numWorkers) // counting semaphore pattern

	for len(queue) > 0 {
		size := len(queue)
		for i := 0; i < size; i++ {
			currentURL := queue[i].currentURL
			pathToCurrentURL := queue[i].pathToCurrentURL
			degree := queue[i].degree

			if urlType == TYPE_PERSON {
				degree += 1
			}

			if currentURL == endURL {
				queue[i].degree = degree
				return &queue[i], nil
			}

			semaphore <- struct{}{} // Acquire a semaphore to control concurrency
			wg.Add(1)

			go func(url string) {
				defer func() {
					<-semaphore // Release the semaphore
					wg.Done()
				}()
				if response, err := fetchMoviebuffData(url, urlType); err == nil {
					mu.Lock()
					visited[currentURL] = true
					mu.Unlock()
					urltype := response.urlType()
					if urltype == "Person" {
						resp := response.(PersonData)
						for _, movie := range resp.Movies {
							path := make([]string, len(pathToCurrentURL))
							copy(path, pathToCurrentURL)
							mu.Lock()
							if !visited[movie.URL] {
								currentURL = movie.URL
								path = append(path, movie.URL)
								q := queueData{
									currentURL:       currentURL,
									degree:           degree,
									pathToCurrentURL: path,
								}
								queue = append(queue, q)
							}
							mu.Unlock()
						}
					} else if urltype == "Movie" {
						resp := response.(MovieData)
						for _, cast := range resp.Cast {
							path := make([]string, len(pathToCurrentURL))
							copy(path, pathToCurrentURL)
							actorURL := cast.URL
							mu.Lock()
							if !visited[actorURL] {
								currentURL = actorURL
								path = append(path, actorURL)
								q := queueData{
									currentURL:       currentURL,
									degree:           degree,
									pathToCurrentURL: path,
								}
								queue = append(queue, q)
							}
							mu.Unlock()
						}
					}
				}
			}(currentURL)
		}

		// Wait for all workers to finish processing
		wg.Wait()
		urlType = 1 - urlType
		queue = queue[size:]
	}

	return nil, fmt.Errorf("no connection found")
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: degrees <person1-url> <person2-url>")
		return
	}

	person1URL := os.Args[1]
	person2URL := os.Args[2]

	result, err := findDegreesOfSeparation(person1URL, person2URL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Degrees of Separation:", result.degree)
	fmt.Println("Path of Separation:")
	fmt.Println(strings.Join(result.pathToCurrentURL, " -> "))
}
