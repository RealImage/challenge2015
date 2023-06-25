package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const (
	baseURL = "http://data.moviebuff.com/"
)

type Person struct {
	URL    string  `json:"url"`
	Name   string  `json:"name"`
	Type   string  `json:"type"`
	Movies []Movie `json:"movies"`
}

type Movie struct {
	URL  string `json:"url"`
	Name string `json:"name"`
	Role string `json:"role"`
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <person1_url> <person2_url>")
		return
	}

	person1URL := os.Args[1]
	person2URL := os.Args[2]

	path, err := findShortestPath(person1URL, person2URL)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Degrees of Separation:", len(path)-1)
	for i := 0; i < len(path)-1; i++ {
		node := path[i]
		nextNode := path[i+1]
		movie, err := getMovie(node.URL)
		fmt.Println(movie)

		if isMovieURL(node.URL) {
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			fmt.Printf("%d. Movie: %s\n", i+1, movie.Name)

			if isActorURL(nextNode.URL) {
				actor, err := getPerson(nextNode.URL)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					return
				}

				printActor(movie, actor)
			} else if isDirectorURL(nextNode.URL) {
				director, err := getPerson(nextNode.URL)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					return
				}

				printDirector(movie, director)
			} else if isSupportingActorURL(nextNode.URL) {
				supportingActor, err := getPerson(nextNode.URL)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					return
				}

				printSupportingActor(movie, supportingActor)
			}
		} else if isActorURL(node.URL) {
			actor, err := getPerson(node.URL)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			printActor(movie, actor)
		} else if isDirectorURL(node.URL) {
			director, err := getPerson(node.URL)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			printDirector(movie, director)
		} else if isSupportingActorURL(node.URL) {
			supportingActor, err := getPerson(node.URL)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			printSupportingActor(movie, supportingActor)
		}
	}
}

func findShortestPath(person1URL, person2URL string) ([]Person, error) {
	queue := make([][]Person, 0)
	visited := make(map[string]bool)

	start := Person{URL: person1URL}
	queue = append(queue, []Person{start})
	fmt.Println(queue)
	visited[start.URL] = true

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]

		node := path[len(path)-1]
		if node.URL == person2URL {
			return path, nil
		}

		neighbors, err := getNeighbors(node.URL)
		if err != nil {
			return nil, err
		}

		for _, neighbor := range neighbors {
			if !visited[neighbor.URL] {
				// newPath := append(path, neighbor)
				// queue = append(queue, newPath)
				visited[neighbor.URL] = true
			}
		}
	}

	return nil, fmt.Errorf("no path found")
}

func getNeighbors(url string) ([]Movie, error) {
	resp, err := http.Get(baseURL + url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var person Person
	err = json.Unmarshal(body, &person)
	if err != nil {
		return nil, err
	}

	return person.Movies, nil
}

func getPerson(url string) (Person, error) {
	resp, err := http.Get(baseURL + url)
	if err != nil {
		return Person{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Person{}, err
	}

	var person Person
	err = json.Unmarshal(body, &person)
	if err != nil {
		return Person{}, err
	}

	return person, nil
}

func getMovie(url string) (Movie, error) {
	resp, err := http.Get(baseURL + url)
	if err != nil {
		return Movie{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Movie{}, err
	}

	var movie Movie
	err = json.Unmarshal(body, &movie)
	if err != nil {
		return Movie{}, err
	}

	return movie, nil
}

func isMovieURL(url string) bool {
	return strings.HasPrefix(url, "movie")
}

func isActorURL(url string) bool {
	return strings.HasPrefix(url, "person/actor")
}

func isDirectorURL(url string) bool {
	return strings.HasPrefix(url, "person/director")
}

func isSupportingActorURL(url string) bool {
	return strings.HasPrefix(url, "person/supporting-actor")
}

func printActor(movie Movie, nextNode Person) {
	for _, m := range nextNode.Movies {
		if m.URL == movie.URL {
			fmt.Printf("Actor: %s\n", nextNode.Name)
			fmt.Printf("Movie: %s\n", movie.Name)
			return
		}
	}
}

func printDirector(movie Movie, nextNode Person) {
	for _, m := range nextNode.Movies {
		if m.URL == movie.URL {
			fmt.Printf("Director: %s\n", nextNode.Name)
			fmt.Printf("Movie: %s\n", movie.Name)
			return
		}
	}
}

func printSupportingActor(movie Movie, nextNode Person) {
	for _, m := range nextNode.Movies {
		if m.URL == movie.URL {
			fmt.Printf("Supporting Actor: %s\n", nextNode.Name)
			fmt.Printf("Movie: %s\n", movie.Name)
			return
		}
	}
}
