package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

func GetActorsData(urls []string) ([]*ActorsData, error) {
	result := make([]*ActorsData, 0)
	for _, url := range urls {
		byteRes, err := hitUrl(url)
		if err != nil {
			return nil, err
		}
		data := &ActorsData{}
		err = json.Unmarshal(byteRes, data)
		if err != nil {
			log.Fatalf("Error  parsing from url. buff : %+s err : %+v", byteRes, err)
			return nil, err
		}
		result = append(result, data)
	}
	return result, nil
}

func hitUrl(url string) ([]byte, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(500)*time.Second)
	defer cancelFunc()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Fatalf("Error in creating request: err : %+v", err)
		return nil, err
	}
	reponse, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatalf("Error in gettting data from url. : err : %+v", err)
		return nil, err
	}
	byteRes, err := io.ReadAll(reponse.Body)
	if err != nil {
		log.Fatalf("Error in reading data from response body. : err : %+v", err)
		return nil, err
	}
	return byteRes, nil
}

func GetMoviesData(urls []string) ([]*MoivesData, error) {

	result := make([]*MoivesData, 0)
	for _, url := range urls {
		byteRes, err := hitUrl(url)
		if err != nil {
			return nil, err
		}
		data := &MoivesData{}
		err = json.Unmarshal(byteRes, data)
		if err != nil {
			log.Fatalf("Error  parsing from url. buff : %+s err : %+v", byteRes, err)
			return nil, err
		}
		result = append(result, data)
	}
	return result, nil
}

func CreateActorDataSet(actorData []*ActorsData, movieData []*MoivesData) map[string]map[string]string {
	var result = make(map[string]map[string]string)
	for _, actor := range actorData {
		GlobalUrlActorMap[actor.URL] = actor.Name
		result[actor.Name] = make(map[string]string)
		for _, movie := range actor.Movies {
			result[actor.Name][movie.Name] = movie.Role
		}
	}

	for _, movie := range movieData {
		for _, cast := range movie.Casts {
			_, ok := result[cast.Name]
			if !ok {
				result[cast.Name] = make(map[string]string)
			}
			result[cast.Name][movie.Name] = cast.Role
		}
	}
	return result
}
func CreateMovieDataSet(actorData []*ActorsData, movieData []*MoivesData) map[string]map[string]string {
	var result = make(map[string]map[string]string)
	for _, movie := range movieData {
		result[movie.Name] = make(map[string]string)
		for _, cast := range movie.Casts {
			GlobalUrlActorMap[cast.URL] = cast.Name
			result[movie.Name][cast.Name] = cast.Role
		}
	}
	for _, actor := range actorData {
		for _, movie := range actor.Movies {
			_, ok := result[movie.Name]
			if !ok {
				result[movie.Name] = make(map[string]string)
			}
			result[movie.Name][actor.Name] = movie.Role
		}
	}
	return result
}

func MakeActorsGraphConnectedWithMovieEdge(movieDataSet, actorDataSet map[string]map[string]string) map[string]map[string]string {
	var graph = make(map[string]map[string]string)
	for actor, setOfMovies := range actorDataSet {
		graph[actor] = make(map[string]string)
		var allTypeOfMovies []string
		for k := range setOfMovies {
			allTypeOfMovies = append(allTypeOfMovies, k)
		}
		allMovies := allTypeOfMovies
		for _, individualMovie := range allMovies {
			var allTypeOfCast []string
			for k := range movieDataSet[individualMovie] {
				allTypeOfCast = append(allTypeOfCast, k)
			}
			for _, cast := range allTypeOfCast {
				graph[actor][cast] = individualMovie
			}
		}
	}
	return graph
}

func FindShortestPath(graph map[string]map[string]string, source, dest string) []MovieNameWithCastName {
	queue := make([]string, 0)
	visited := make(map[string]bool)
	parents := make(map[string]ActorMovieRelation)
	if _, ok := GlobalUrlActorMap[source]; !ok {
		log.Fatal("source not exist in data set")
		return nil
	} else if _, ok := GlobalUrlActorMap[dest]; !ok {
		log.Fatal("dest not exist in data set")
		return nil
	}
	queue = append(queue, GlobalUrlActorMap[source])
	parents[GlobalUrlActorMap[source]] = ActorMovieRelation{Actor: GlobalUrlActorMap[source], commonMovie: ""}
	i := 0
	for i < len(queue) {
		currentNode := queue[i]
		visited[currentNode] = true
		var allKeys []string
		for k := range graph[currentNode] {
			allKeys = append(allKeys, k)
		}
		for _, key := range allKeys {
			ok := visited[key]
			if !ok {
				queue = append(queue, key)
				visited[key] = true
				parents[key] = ActorMovieRelation{Actor: currentNode, commonMovie: graph[currentNode][key]}
			}
		}
		i++
	}
	temp := GlobalUrlActorMap[dest]
	var result []MovieNameWithCastName
	for temp != GlobalUrlActorMap[source] {
		p := parents[temp]
		result = append(result, MovieNameWithCastName{Cast1: p.Actor, Cast2: temp, Movie: p.commonMovie})
		temp = p.Actor
	}
	// Backtrack
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return result
}
