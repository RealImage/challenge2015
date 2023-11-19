package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

const moviebuffURL = "https://data.moviebuff.com/%s"

type Movie struct {
	URL  string `json:"url"`
	Name string `json:"name"`
	Type string `json:"type"`
	Cast []Role `json:"cast"`
	Crew []Role `json:"crew"`
}

type Role struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Role string `json:"role"`
}

type Person struct {
	URL    string `json:"url"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Movies []Role `json:"movies"`
}

func fetchPersonData(url string) (*Person, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// fmt.Println(url)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data from %s: %s", url, resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var person Person
	err = json.Unmarshal(body, &person)
	if err != nil {
		return nil, err
	}

	return &person, nil
}

func fetchMovieData(url string) (*Movie, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// fmt.Println(url)

	if resp.StatusCode != http.StatusOK {
		// return nil, fmt.Errorf("failed to fetch data from %s: %s", url, resp.Status)
		fmt.Errorf("failed to fetch info for ", url)
		return nil, nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var movie Movie
	err = json.Unmarshal(body, &movie)
	if err != nil {
		return nil, err
	}

	return &movie, nil
}

// Movie[the-great-gatsby]=[]cast,crew
var MovieMap map[string]Movie

// Person[amitabh-bachchan]=[]movies
var PersonMap map[string][]Movie

var wg sync.WaitGroup

// var CrewMap map[string][]Role

var visited map[string]bool
var queue [][]*Person

func init() {
	visited = make(map[string]bool)
	queue = make([][]*Person, 0)
	MovieMap = make(map[string]Movie)

	PersonMap = make(map[string][]Movie)
}

func findDegreesOfSeparation(start, end string) ([]*Person, error) {

	var mu sync.Mutex

	// var result []*Person

	startPerson, err := fetchPersonData(fmt.Sprintf(moviebuffURL, start))
	if err != nil {
		return nil, err
	}

	// endPerson, err := fetchPersonData(fmt.Sprintf(moviebuffURL, end))
	// if err != nil {
	// 	return nil, err
	// }

	for _, movie := range startPerson.Movies {
		thisMovie, err := fetchMovieData(fmt.Sprintf(moviebuffURL, movie.URL))
		if err != nil {
			return nil, err
		}
		if thisMovie == nil {
			continue
		}

		if _, seen := MovieMap[movie.URL]; !seen {
			MovieMap[movie.URL] = *thisMovie
		}
		if _, seen := PersonMap[startPerson.URL]; !seen {
			PersonMap[startPerson.URL] = append(PersonMap[startPerson.URL], *thisMovie)
		}
		for _, crew := range thisMovie.Crew {
			if _, seen := PersonMap[crew.URL]; !seen {
				PersonMap[startPerson.URL] = append(PersonMap[startPerson.URL], *thisMovie)
			}

		}
	}

	fmt.Println("After calculating first person")
	fmt.Println("length of moviemap -> ", len(MovieMap))
	fmt.Println("length of personmap -> ", len(PersonMap))

	wg.Add(5)
	go func() {
		mu.Lock()
		movies, seen := PersonMap[end]
		if seen {
			defer wg.Done()
			fmt.Println("end person is found")
			fmt.Println("movie list -> ", movies)
		}

		for _, movie := range MovieMap {
			for _, person := range movie.Cast {
				if _, seen := PersonMap[person.URL]; seen {
					continue
				}
				thisPerson, err := fetchPersonData(fmt.Sprintf(moviebuffURL, person.URL))
				if thisPerson == nil {
					continue
				}
				if err != nil {
					fmt.Println("failed to fetch person - >", thisPerson.Name)
				}

				PersonMap[thisPerson.URL] = append(PersonMap[thisPerson.URL], movie)
			}
		}
		mu.Unlock()
		time.Sleep(1 * time.Minute)
	}()

	go func() {
		mu.Lock()
		movies, seen := PersonMap[end]
		if seen {
			defer wg.Done()
			fmt.Println("end person is found")
			fmt.Println("movie list -> ", movies)
		}

		for _, movie := range MovieMap {
			for _, person := range movie.Crew {
				if _, seen := PersonMap[person.URL]; seen {
					continue
				}
				thisPerson, err := fetchPersonData(fmt.Sprintf(moviebuffURL, person.URL))
				if thisPerson == nil {
					continue
				}
				if err != nil {
					fmt.Println("failed to fetch person - >", thisPerson.Name)
				}
				PersonMap[thisPerson.URL] = append(PersonMap[thisPerson.URL], movie)
			}
		}
		mu.Unlock()
		time.Sleep(1 * time.Minute)
	}()

	go func() {
		mu.Lock()
		movies, seen := PersonMap[end]
		if seen {
			defer wg.Done()
			fmt.Println("end person is found")
			fmt.Println("movie list -> ", movies)
		}

		for _, movies := range PersonMap {
			for _, movie := range movies {
				if _, seen := MovieMap[movie.URL]; seen {
					continue
				}
				thisMovie, err := fetchMovieData(fmt.Sprintf(moviebuffURL, movie.URL))
				if err != nil {
					fmt.Println("failed to fetch movie info -> ", movie.URL)
				}
				if thisMovie == nil {
					continue
				}

				MovieMap[movie.URL] = *thisMovie
				if _, seen := PersonMap[startPerson.URL]; !seen {
					PersonMap[startPerson.URL] = append(PersonMap[startPerson.URL], *thisMovie)
				}
				for _, crew := range thisMovie.Crew {
					if _, seen := PersonMap[crew.URL]; !seen {
						PersonMap[startPerson.URL] = append(PersonMap[startPerson.URL], *thisMovie)
					}

				}
			}
		}
		mu.Unlock()
		time.Sleep(1 * time.Minute)
	}()

	go func() {

		for {
			mu.Lock()
			movies, seen := PersonMap[end]
			if seen {
				defer wg.Done()
				fmt.Println("end person is found")
				fmt.Println("movie list -> ", movies)
			}
			fmt.Println("length of moviemap -> ", len(MovieMap))
			fmt.Println("length of personmap -> ", len(PersonMap))
			mu.Unlock()
			time.Sleep(30 * time.Second)
		}

	}()

	// write to file
	go func() {
		for {
			mu.Lock()
			personBytes := new(bytes.Buffer)
			person := gob.NewEncoder(personBytes)

			// Encoding the map
			err := person.Encode(PersonMap)
			if err != nil {
				panic(err)
			}

			movieBytes := new(bytes.Buffer)
			movie := gob.NewEncoder(movieBytes)

			// Encoding the map
			err = person.Encode(movie)
			if err != nil {
				panic(err)
			}

			if err := ioutil.WriteFile("personmap", personBytes.Bytes(), 777); err != nil {
				fmt.Println("failed to write to file")
			}

			if err := ioutil.WriteFile("moviemap", movieBytes.Bytes(), 777); err != nil {
				fmt.Println("failed to write to file")
			}
			mu.Unlock()
			time.Sleep(2 * time.Minute)
		}
	}()

	// match - find if end person is present in the person map
	go func() {
		for {
			mu.Lock()
			movies, seen := PersonMap[end]
			if seen {
				defer wg.Done()
				fmt.Println("end person is found")
				fmt.Println("movie list -> ", movies)
				break
			}
			mu.Unlock()
			time.Sleep(5 * time.Minute)
		}

	}()

	// queue = append(queue, []*Person{startPerson})

	// for len(queue) > 0 {
	// 	currentPath := queue[0]
	// 	queue = queue[1:]

	// 	currentPerson := currentPath[len(currentPath)-1]

	// 	if currentPerson.Name == end {
	// 		return currentPath, nil
	// 	}

	// 	if visited[currentPerson.Name] {
	// 		continue
	// 	}

	// 	visited[currentPerson.Name] = true

	// 	for _, movie := range currentPerson.Movies {
	// 		movieURL := fmt.Sprintf(moviebuffURL, movie.URL)
	// 		movieData, err := fetchMovieData(movieURL)
	// 		if err != nil {
	// 			return nil, err
	// 		}

	// 		if movieData == nil {
	// 			continue
	// 		}

	// 		var relatedPeople []*Person
	// 		for _, cast := range movieData.Cast {
	// 			person, err := fetchPersonData(fmt.Sprintf(moviebuffURL, cast.URL))
	// 			if err != nil {
	// 				return nil, err
	// 			}
	// 			relatedPeople = append(relatedPeople, person)
	// 		}

	// 		// relatedPeople := []*Person{movieData}

	// 		for _, relatedPerson := range relatedPeople {
	// 			if !visited[relatedPerson.Name] {
	// 				newPath := append(currentPath, relatedPerson)
	// 				queue = append(queue, newPath)
	// 			}
	// 		}
	// 	}

	// }

	// return result, fmt.Errorf("no connection found between %s and %s", start, end)

	return nil, nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: degrees <start_person_url> <end_person_url>")
		os.Exit(1)
	}

	startPerson := os.Args[1]
	endPerson := os.Args[2]

	result, err := findDegreesOfSeparation(startPerson, endPerson)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	wg.Wait()

	fmt.Printf("Degrees of Separation: %d\n", len(result)-1)
	for i, person := range result {
		fmt.Printf("%d. Movie: %s\n   %s: %s\n", i+1, person.Movies, getPersonType(person), person.Name)
	}
}

func getPersonType(person *Person) string {
	if len(person.Movies) > 0 {
		return "Supporting Actor"
	}
	return "Actor"
}
