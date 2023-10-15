package movie

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"myproject/challenge2015/dtos"
	"net/http"
	"strconv"
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
	personURL string
	associate []*PersonNetwork
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

	// movieQueue := &Queue{}
	networkQueue := &Queue{}
	networkQueue.Enqueue(person1Node)
	queueLen := 0
	for !networkQueue.IsEmpty() && !isPerson2Found {
		if queueLen == 0 && !isPerson2Found {
			degree++
			queueLen = networkQueue.Size()
		}
		person := networkQueue.Dequeue()
		queueLen--
		if root == nil {
			root = person
		}

		personDetails, err := m.GetPersonDetailsByURL(person.personURL)
		if err != nil {
			return degree, err
		}
		for _, movie := range personDetails.MoviesAndRoles {
			//if _, ok := movieDiscovered[movie.URL]; !ok { // if we ignore this movie for current person then some edge for current emp may missed.
			movieDetails, err := m.GetMovieDetailsByURL(movie.URL) //TODO: we can store the this movie details if already fetched
			if err != nil {
				return degree, err
			}

			//Process cast
			for _, cast := range movieDetails.Cast {
				if cast.URL == person2 {
					isPerson2Found = true
				}
				if _, ok := personDiscovered[cast.URL]; !ok {

					castPersonNode := &PersonNetwork{
						personURL: cast.URL,
						associate: []*PersonNetwork{},
					}
					personDiscovered[cast.URL] = castPersonNode
					person.associate = append(person.associate, castPersonNode)
					networkQueue.Enqueue(castPersonNode)

				} else {
					if person.personURL != cast.URL {
						person.associate = append(person.associate, personDiscovered[cast.URL])
					}
				}
			}

			//Process crew
			for _, crew := range movieDetails.Crew {
				if crew.URL == person2 {
					isPerson2Found = true
				}
				if _, ok := personDiscovered[crew.URL]; !ok {

					crewPersonNode := &PersonNetwork{
						personURL: crew.URL,
						associate: []*PersonNetwork{},
					}
					personDiscovered[crew.URL] = crewPersonNode
					person.associate = append(person.associate, crewPersonNode)
					networkQueue.Enqueue(crewPersonNode)

				} else {
					if person.personURL != crew.URL {
						person.associate = append(person.associate, personDiscovered[crew.URL])
					}
				}
			}
			//}
		}

	}

	return degree, nil
}

func (m *MovieService) GetPersonDetailsByURL(personUrl string) (dtos.ActorDetails, error) {

	response := dtos.ActorDetails{}

	for i := 0; i < 4; i++ { // retry 3 time

		resp, err := http.Get("https://data.moviebuff.com/" + personUrl)

		if resp.StatusCode == 429 {
			waitTime, _ := strconv.Atoi(resp.Header.Get("retry-after"))
			time.Sleep(time.Duration(waitTime) * time.Second)

		} else if err != nil {
			return response, err
		}

		if resp.StatusCode == 200 {

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

	for i := 0; i < 4; i++ { // retry 3 time

		resp, err := http.Get("https://data.moviebuff.com/" + movieNameUrl)

		if resp.StatusCode == 429 {
			waitTime, _ := strconv.Atoi(resp.Header.Get("retry-after"))
			time.Sleep(time.Duration(waitTime) * time.Second)

		} else if err != nil {
			return response, err
		}

		if resp.StatusCode == 200 {

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

//TODO check for the case where the get fucntions returns empty.
