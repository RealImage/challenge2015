package services

import (
	domainmodels "challenge2015/internal/domain/models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	FirstActor  string
	SecondActor string
	// actor->movie->role
	actorsDataMap = map[string]map[string]string{}
	// movie->cast
	movieDataMap             = map[string][]domainmodels.Data{}
	visitedActorsMap         = map[string]bool{}
	firstActorMoviesVisited  = map[string]bool{}
	secondActorMoviesVisited = map[string]bool{}
	ActorListForFirst        = []string{}
	ActorListForSecond       = []string{}
	moviesListForFirstActor  = []string{}
	moviesListForSecondActor = []string{}
)

func SmallestDegreeOfSeparation() error {

	commonMovie, err := tryFindingCommonMovie()
	if err != nil {
		return fmt.Errorf("no connection found between the actors")
	}
	if commonMovie != "" {
		findDegreeOfSeparation(FirstActor, SecondActor)
		return nil
	}

	err = appendCastMemberToActorsList()
	if err != nil {
		return fmt.Errorf("no connection found")
	}

	return SmallestDegreeOfSeparation()
}

// finds common movie between two actors and returns if any.
func tryFindingCommonMovie() (string, error) {
	if len(ActorListForFirst) == 0 || len(ActorListForSecond) == 0 {
		return "", errors.New("try finding common movie: no actors found")
	}

	for len(ActorListForFirst) > 0 && len(ActorListForSecond) > 0 {

		firstActor := ActorListForFirst[0]
		ActorListForFirst = ActorListForFirst[1:]
		actorsDataMap[firstActor] = map[string]string{}

		secondActor := ActorListForSecond[0]
		ActorListForSecond = ActorListForSecond[1:]
		actorsDataMap[secondActor] = map[string]string{}

		moviesDataForFirstActor, err := getActorsData(firstActor)
		if err != nil {
			return "", nil
		}
		moviesDataForSecondActor, err := getActorsData(secondActor)
		if err != nil {
			return "", nil
		}

		for _, movie := range moviesDataForFirstActor.Movies {
			actorsDataMap[firstActor][movie.Url] = movie.Role
			if _, ok := firstActorMoviesVisited[movie.Url]; !ok {
				firstActorMoviesVisited[movie.Url] = true
				moviesListForFirstActor = append(moviesListForFirstActor, movie.Url)
			}
		}

		for _, movie := range moviesDataForSecondActor.Movies {
			actorsDataMap[secondActor][movie.Url] = movie.Role
			if _, ok := secondActorMoviesVisited[movie.Url]; !ok {
				secondActorMoviesVisited[movie.Url] = true
				moviesListForSecondActor = append(moviesListForSecondActor, movie.Url)
			}
		}

		commonMovie := searchCommonMovie(firstActorMoviesVisited, secondActorMoviesVisited)
		if commonMovie != "" {
			return commonMovie, nil
		}
	}

	return "", nil
}

// appends all cast members to actors list for first and second actor.
func appendCastMemberToActorsList() error {
	if len(moviesListForFirstActor) == 0 || len(moviesListForSecondActor) == 0 {
		return errors.New("try finding common cast member: empty movie lists")
	}

	for len(moviesListForFirstActor) > 0 && len(moviesListForSecondActor) > 0 {

		firstActorMovie := moviesListForFirstActor[0]
		moviesListForFirstActor = moviesListForFirstActor[1:]

		secondActorMovie := moviesListForSecondActor[0]
		moviesListForSecondActor = moviesListForSecondActor[1:]

		moviesDataForFirstActor, err := getMovieData(firstActorMovie)
		if err != nil {
			return nil
		}
		moviesDataForSecondActor, err := getMovieData(secondActorMovie)
		if err != nil {
			return nil
		}

		firstActorsCast := []domainmodels.Data{}
		firstActorsCast = append(firstActorsCast, moviesDataForFirstActor.Cast...)
		firstActorsCast = append(firstActorsCast, moviesDataForFirstActor.Crew...)
		movieDataMap[firstActorMovie] = append(movieDataMap[firstActorMovie], firstActorsCast...)

		for _, cast := range firstActorsCast {
			if _, ok := visitedActorsMap[cast.Name]; !ok {
				visitedActorsMap[cast.Name] = true
				ActorListForFirst = append(ActorListForFirst, cast.Url)
			}
		}

		secondActorsCast := []domainmodels.Data{}
		secondActorsCast = append(secondActorsCast, moviesDataForSecondActor.Cast...)
		secondActorsCast = append(secondActorsCast, moviesDataForSecondActor.Crew...)
		movieDataMap[secondActorMovie] = append(movieDataMap[secondActorMovie], secondActorsCast...)

		for _, cast := range secondActorsCast {
			if _, ok := visitedActorsMap[cast.Name]; !ok {
				visitedActorsMap[cast.Name] = true
				ActorListForSecond = append(ActorListForSecond, cast.Url)
			}
		}

	}

	return nil
}

type QueueItem struct {
	actor        string
	level        int
	commonActors []string
}

// finds degree of separation between the actors and returns.
func findDegreeOfSeparation(firstActor, secondActor string) {
	visited := make(map[string]bool)
	queue := make([]QueueItem, 0)

	queue = append(queue, QueueItem{firstActor, 0, []string{firstActor}})
	visited[firstActor] = true

	for len(queue) > 0 {
		currentItem := queue[0]
		queue = queue[1:]

		currentActor := currentItem.actor
		currentLevel := currentItem.level
		currentcommonActors := currentItem.commonActors

		if currentActor == secondActor {
			print(currentcommonActors)
			return
		}

		for movie := range actorsDataMap[currentActor] {
			for _, nextActor := range movieDataMap[movie] {
				if !visited[nextActor.Name] {
					visited[nextActor.Name] = true
					nextcommonActors := append(currentcommonActors, nextActor.Url)
					queue = append(queue, QueueItem{nextActor.Url, currentLevel + 1, nextcommonActors})
				}
			}
		}
	}
}

// output.
func print(commonActors []string) {
	fmt.Println("Degrees of separation: ", len(commonActors)-1)
	for i := 0; i < len(commonActors)-1; i++ {
		fmt.Println()
		movie, role1, role2 := findCommonMovie(actorsDataMap[commonActors[i]], actorsDataMap[commonActors[i+1]])
		fmt.Printf("%d. Movie: %s\n", i+1, movie)
		fmt.Printf("%s: %s\n", role1, commonActors[i])
		fmt.Printf("%s: %s\n", role2, commonActors[i+1])
	}
}

// returns common movie from first and second actors movies visited map.
func searchCommonMovie(firstActorMoviesVisited, secondActorMoviesVisited map[string]bool) string {
	for k := range firstActorMoviesVisited {
		if _, ok := secondActorMoviesVisited[k]; ok {
			if len(movieDataMap[k]) > 1 {
				return k
			}
		}
	}
	return ""
}

func findCommonMovie(mv1, mv2 map[string]string) (string, string, string) {
	for m1, r1 := range mv1 {
		for m2, r2 := range mv2 {
			if m1 == m2 {
				return m1, r1, r2
			}
		}
	}
	return "", "", ""
}

func getActorsData(actor string) (domainmodels.ActorData, error) {
	data, err := getData(actor)
	if err != nil {
		return domainmodels.ActorData{}, err
	}
	if data.Type != "Person" {
		return domainmodels.ActorData{}, errors.New("get actors data: invalid data")
	}
	return domainmodels.ActorData{
		Url:    data.Url,
		Type:   data.Type,
		Name:   data.Name,
		Movies: data.Movies,
	}, nil
}

func getMovieData(movie string) (domainmodels.MovieData, error) {
	data, err := getData(movie)
	if err != nil {
		return domainmodels.MovieData{}, err
	}
	if data.Type != "Movie" {
		return domainmodels.MovieData{}, errors.New("get movie data: invalid data")
	}
	return domainmodels.MovieData{
		Url:  data.Url,
		Type: data.Type,
		Name: data.Name,
		Cast: data.Cast,
		Crew: data.Crew,
	}, nil
}

func getData(input string) (domainmodels.Data, error) {
	url := fmt.Sprintf("http://data.moviebuff.com/%s", input)
	// Make the HTTP request
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return domainmodels.Data{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return domainmodels.Data{}, errors.New("Invalid")
	}

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return domainmodels.Data{}, err
	}

	var data domainmodels.Data
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return domainmodels.Data{}, err
	}

	return data, nil
}
