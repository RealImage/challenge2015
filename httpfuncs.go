package main

import (
	"fmt"
	"os"
)

const moviebuff = "http://data.moviebuff.com/"

var (
	movieList map[string][]string
	seen      map[string]bool
	actors    []string
	degrees   int
)

func ErrHandle(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Print("Usage Example : degrees vn-mayekar magie-mathur")
	}
	seen = make(map[string]bool)
	buildTree(os.Args[1], os.Args[2])
}

func buildTree(argument, destination string) {
	degrees++
	actors := []string{}
	newActors := actors
	buildActors(argument, destination)
	for _, actor := range newActors{
		buildTree(actor, destination)
	}

}

func buildActors(argument, destination string) {
	url := moviebuff + argument
	json, err := getData(url)
	defer ErrHandle(err)

	for _, movie := range json.Movies {
		if isSeen(movie.Url) {
			buildActors(movie.Url, destination)
		}
	}

	for _, cast := range json.Cast {
		if isSeen(cast.Url) {
			actors = append(actors, cast.Url)
			if cast.Url == destination {
				fmt.Println("DONE --> ", cast.Url, degrees)
				os.Exit(1)
			}
		}
	}
}

func isSeen(in string) bool {
	if !seen[in] {
		seen[in] = true
		return true
	} else {
		return false
	}
}
