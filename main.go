package main

import (
	_ "embed"
	"log"
)

// In this solution I have used maps extensively because of the fast lookup time than slices.
// maps makes following queries vert fast :
// 1. Is Actor A in cast of movie B ?
// 2. Is Movie M contains Actor A as cast ?
// 3. Wether Actor A and Actor B are in the cast of movie M ?

// I have also created a custom data type PairSet which is struct with only one field that is a map[string]string.
// Set is implemented through maps due to followign reason:
// 1. Fast lookup
// 2. Automatic duplicate removal by Golang.

// My Approach :
// 1. Make a graph using maps where keys are actor names and values are themselves a map of key actor names and value movie.
// example : {"Actor A" : map[string]string{"Actor B" : "Movie M"}} --> This represent Actor A is related to Actor B using through M.
// 2. Find the shortest path between two nodes using BFS.

//go:embed urls.json
var bytes []byte

const URL_FILE = "" // Just for testing using go run *.go command

// BFS: Performs the Breadth-first-search on the graph to find the shortest path between two
// actors using common movie as an edge.
func BFS(graph map[string]*PairSet, source, dest string) []MovieNameCastName {
	// Temporary struct just to represent a parent in the graph
	type parent struct {
		parentActor string
		commonMovie string
	}
	queue := make([]string, 0)
	visited := make(map[string]bool)
	parents := make(map[string]parent)

	queue = append(queue, source)
	parents[source] = parent{parentActor: source, commonMovie: ""}

	var i int
	for i < len(queue) {
		elem := queue[i]
		visited[elem] = true
		for _, v := range graph[elem].GetAllKeys() {
			ok := visited[v]
			if !ok {
				queue = append(queue, v)
				visited[v] = true
				parents[v] = parent{parentActor: elem, commonMovie: graph[elem].Get(v)}
			}
		}
		i++
	}
	temp := dest
	var result []MovieNameCastName
	for temp != source {
		p := parents[temp]
		result = append(result, MovieNameCastName{Cast1: p.parentActor, Cast2: temp, Movie: p.commonMovie})
		temp = p.parentActor
	}
	return result
}

func main() {
	nodes := GetCmdLineArgs()
	source := nodes[0]
	destination := nodes[1]

	urls, err := ParseURLsFile(URL_FILE, bytes)
	if err != nil {
		log.Fatal(err)
	}

	actors, err := FetchActorsData(urls.Actor)
	if err != nil {
		log.Fatal(err)
	}

	movies, err := FetchMoviesData(urls.Movies)
	if err != nil {
		log.Fatal(err)
	}

	// {movie : {actor:role, actor:role, actor:role}, movie : {actor:role, actor:role, actor:role}}
	movieMap := MakeMovieMap(actors, movies)

	// {actor : {movie:role, movie:role, movie:role}, actor : {movie:role, movie:role, movie:role}}
	actorMap := MakeActorMap(actors, movies)

	// Adjacency List
	// {actor : {actor:movie, actor:movie, actor:movie}, actor : {actor:movie, actor:movie, actor:movie}}
	graph := MakeAdjecencyList(movieMap, actorMap)

	shortestPath := BFS(graph, source, destination)
	shortestPath = Reverse(shortestPath)
	PrintResult(shortestPath, actorMap, movieMap, source, destination)
}
