package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/spf13/pflag"
)

func main() {
	var source, destination string
	pflag.StringVarP(&source, "source", "s", "", "actor source")
	pflag.StringVarP(&destination, "destination", "d", "", "actor destination")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	if source == "" || destination == "" {
		log.Fatal("use this: go run *.go -s <actorname> -d <actorname>")
		os.Exit(0)
	}
	actors, err := GetActorsData(ActorsUrl)
	if err != nil {
		log.Fatal(err)
	}
	movies, err := GetMoviesData(MoviesUrl)
	if err != nil {
		log.Fatal(err)
	}
	movieData := CreateMovieDataSet(actors, movies)
	actorData := CreateActorDataSet(actors, movies)
	graph := MakeActorsGraphConnectedWithMovieEdge(movieData, actorData)
	shortestPath := FindShortestPath(graph, source, destination)
	// Print Result
	fmt.Println("Degrees of Separation: ", len(shortestPath))
	fmt.Println("----------------------------------------------")
	for i, v := range shortestPath {
		fmt.Printf("%d. %s\n", i+1, v.Movie)
		fmt.Printf("%s: %s\n", actorData[v.Cast1][v.Movie], v.Cast1)
		fmt.Printf("%s: %s\n", actorData[v.Cast2][v.Movie], v.Cast2)
		fmt.Println("----------------------------------------------")
	}

}
