package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const baseURL = "http://data.moviebuff.com/"

type Person struct {
	URL    string   `json:"url"`
	Type   string   `json:"type"`
	Name   string   `json:"name"`
	Movies []Detail `json:"movies"`
}
type Detail struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Role string `json:"role"`
}
type Movie struct {
	URL  string   `json:"url"`
	Type string   `json:"type"`
	Name string   `json:"name"`
	Cast []Detail `json:"cast"`
	Crew []Detail `json:"crew"`
}

func fetchJSON(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching data from %s: %s", url, err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error fetching data from %s: Status code %d", url, resp.StatusCode)
		return fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		log.Printf("Error fetching data from %s: Unexpected content type: %s", url, contentType)
		return fmt.Errorf("Unexpected content type: %s", contentType)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body from %s: %s", url, err)
		return err
	}

	return json.Unmarshal(body, target)
}

func findDegreesOfSeparation(actor1, actor2 string) (int, []string) {
	queue := [][]string{{actor1}}
	visited := make(map[string]bool)
	visited[actor1] = true

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]

		currentActor := path[len(path)-1]

		if currentActor == actor2 {
			return len(path) - 1, path
		}

		// Fetch movies for the current actor
		actorURL := baseURL + currentActor
		actorData := new(Person)
		if err := fetchJSON(actorURL, actorData); err != nil {
			log.Printf("Error fetching data for %s: %s", currentActor, err)
			continue
		}
		//log.Println("Actor data~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~", actorURL, actorData)
		for _, movie := range actorData.Movies {
			// Fetch cast and crew for the movie
			movieURL := baseURL + movie.URL
			movieData := new(Movie)
			if err := fetchJSON(movieURL, movieData); err != nil {
				log.Printf("Error fetching data for %s: %s", movie.URL, err)
				continue
			}

			for _, person := range append(movieData.Cast, movieData.Crew...) {
				if !visited[person.URL] {
					visited[person.URL] = true
					newPath := append(path, person.URL)
					queue = append(queue, newPath)
				}
			}
		}
	}

	return -1, nil
}

func printDegreesOfSeparation(degrees int, path []string) {
	fmt.Printf("Degrees of Separation: %d\n\n", degrees)

	for i := 1; i < len(path); i++ {
		fmt.Printf("%d. Movie: %s\n", i, path[i-1])
		fmt.Printf("   %s: %s\n", strings.Title(path[i]), path[i])
	}

	fmt.Printf("%d. Movie: %s\n", len(path), path[len(path)-1])
}

func main() {
	logFile, err := os.OpenFile("info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error opening log file: ", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	if len(os.Args) != 4 {
		log.Println("Usage: go run main.go <actor1> <actor2>")
		os.Exit(1)
	}

	actor1 := os.Args[1]
	actor2 := os.Args[2]

	degrees, path := findDegreesOfSeparation(actor1, actor2)

	if degrees == -1 {
		log.Printf("No connection found between %s and %s", actor1, actor2)
		fmt.Println("No connection found between", actor1, "and", actor2)
	} else {
		log.Printf("Degrees of Separation between %s and %s: %d", actor1, actor2, degrees)
		printDegreesOfSeparation(degrees, path)
	}
}
