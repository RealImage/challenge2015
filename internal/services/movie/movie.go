package movie

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"myproject/challenge2015/dtos"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type MovieService struct{}

func New() *MovieService {
	return &MovieService{}
}

type IMovie interface {
	GetMinimumDegreeOfSeperation(person1, person2 string) (int, error)
}

// Node represents a node in the graph.
type PersonNetwork struct {
	personURL     string
	personRole    []string
	associateRole []string
	movieName     []string
	associate     []*PersonNetwork
}

func (n *PersonNetwork) Add_Edge(node *PersonNetwork) {
	n.associate = append(n.associate, node)
}

func (m *MovieService) GetMinimumDegreeOfSeperation(person1, person2 string) (int, error) {

	degree := 0
	isPerson2Found := false
	personDiscovered := make(map[string]*PersonNetwork)
	// movieDiscovered := make(map[string]*PersonNetwork)

	var root *PersonNetwork

	person1Node := &PersonNetwork{
		personURL: person1,
		associate: []*PersonNetwork{},
	}
	personDiscovered[person1Node.personURL] = person1Node

	if root == nil {
		root = person1Node
	}

	networkQueue := &Queue{}
	networkQueue.Enqueue(person1Node)
	lock := sync.Mutex{}

	for !networkQueue.IsEmpty() && !isPerson2Found {
		if !isPerson2Found {
			degree++
		}
		wtgrp1 := sync.WaitGroup{}
		networkQueueInner := &Queue{} // networkQueueInner.Enqueue() update this inside below loop

		for {

			person := networkQueue.Dequeue()
			if person == nil {
				break
			}

			wtgrp1.Add(1)

			go func(personURL string) {
				defer wtgrp1.Done()

				personDetails, err := m.GetPersonDetailsByURL(personURL)
				if err != nil {
					log.Printf("GetPersonDetailsByURL: %s error: %s", personURL, err.Error())
				}
				for _, movie := range personDetails.MoviesAndRoles {

					movieDetails, err := m.GetMovieDetailsByURL(movie.URL) //TODO: make get movie detail call concurrent using gorutine
					if err != nil {
						log.Printf("GetMovieDetailsByURL: %s error: %s", movie.URL, err.Error())
					}

					//find person role in movie //from both cast and crew.
					personRole := ""
					for _, cast := range movieDetails.Cast {
						if person.personURL != cast.URL {
							personRole = cast.Role
						}
					}
					if len(personRole) == 0 {
						for _, crew := range movieDetails.Crew {
							if person.personURL != crew.URL {
								personRole = crew.Role
							}
						}
					}

					//Process cast
					for _, cast := range movieDetails.Cast {
						if cast.URL == person2 {
							isPerson2Found = true
						}

						lock.Lock()
						if _, ok := personDiscovered[cast.URL]; !ok {

							castPersonNode := &PersonNetwork{
								personURL: cast.URL,
								associate: []*PersonNetwork{},
							}
							personDiscovered[cast.URL] = castPersonNode
							person.associate = append(person.associate, castPersonNode)
							networkQueueInner.Enqueue(castPersonNode)

						} else {
							if person.personURL != cast.URL {
								person.associate = append(person.associate, personDiscovered[cast.URL])
							}
						}
						person.movieName = append(person.movieName, movieDetails.Name)
						person.personRole = append(person.personRole, personRole)
						person.associateRole = append(person.associateRole, cast.Role)

						//TODO: setting a reverse link

						lock.Unlock()
					}

					//Process crew
					for _, crew := range movieDetails.Crew {
						if crew.URL == person2 {
							isPerson2Found = true
						}
						lock.Lock()
						if _, ok := personDiscovered[crew.URL]; !ok {

							crewPersonNode := &PersonNetwork{
								personURL: crew.URL,
								associate: []*PersonNetwork{},
							}
							personDiscovered[crew.URL] = crewPersonNode
							person.associate = append(person.associate, crewPersonNode)
							networkQueueInner.Enqueue(crewPersonNode)

						} else {
							if person.personURL != crew.URL {
								person.associate = append(person.associate, personDiscovered[crew.URL])
							}
						}
						person.movieName = append(person.movieName, movieDetails.Name)
						person.personRole = append(person.personRole, personRole)
						person.associateRole = append(person.associateRole, crew.Role)

						//TODO: setting a reverse link

						lock.Unlock()
					}
					//}
				}
			}(person.personURL)

		}
		wtgrp1.Wait()
		networkQueue = networkQueueInner
	}

	return degree, nil
}

func (m *MovieService) GetPersonDetailsByURL(personUrl string) (dtos.ActorDetails, error) {

	response := dtos.ActorDetails{}

	for i := 0; i < 10; i++ { // retry 10 time

		resp, err := http.Get("https://data.moviebuff.com/" + personUrl)

		if resp != nil && resp.StatusCode == 429 {
			waitTime, _ := strconv.Atoi(resp.Header.Get("retry-after"))
			time.Sleep(time.Duration(waitTime) * time.Second)

		} else if err != nil && i < 10 {
			log.Printf("GetPersonDetailsByURL sleep : %s", err.Error())
			time.Sleep(time.Duration(60) * time.Second)
		} else if err != nil {
			log.Printf("GetPersonDetailsByURL error: %s ", err.Error())
			return response, err
		}

		if resp != nil && resp.StatusCode == 200 {

			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			err = json.Unmarshal(body, &response)
			if err != nil {
				log.Printf("Reading body failed: %s", err)
				return response, err
			}
			break
		}

	}

	return response, nil
}

func (m *MovieService) GetMovieDetailsByURL(movieNameUrl string) (dtos.MovieDetail, error) {

	response := dtos.MovieDetail{}

	for i := 0; i < 10; i++ { // retry 10 time

		resp, err := http.Get("https://data.moviebuff.com/" + movieNameUrl)

		if resp != nil && resp.StatusCode == 429 {
			waitTime, _ := strconv.Atoi(resp.Header.Get("retry-after"))
			time.Sleep(time.Duration(waitTime) * time.Second)

		} else if err != nil && i < 10 {
			log.Printf("GetMovieDetailsByURL sleep : %s", err.Error())
			time.Sleep(time.Duration(60) * time.Second)
		} else if err != nil {
			log.Printf("GetMovieDetailsByURL error: %s", err.Error())
			return response, err
		}

		if resp != nil && resp.StatusCode == 200 {

			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			err = json.Unmarshal(body, &response)
			if err != nil {
				log.Printf("Reading body failed: %s", err)
				return response, err
			}
			break
		}

	}

	return response, nil
}
