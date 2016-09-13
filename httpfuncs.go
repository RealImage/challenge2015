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
	buildTree(os.Args[1], os.Args[2], os.Args[1], true)
}

func buildTree(argument, destination, parent string, getMovie bool) {
	actors = []string{}
	buildActors(argument, destination, parent)
	fmt.Println(actors)
}

func buildActors(argument, destination, parent string) {
	url := moviebuff + argument
	json, err := getData(url)
	defer ErrHandle(err)

	for _, movie := range json.Movies {
		if isSeen(movie.Url) {
			buildActors(movie.Url, destination, argument)
		}
	}

	for _, cast := range json.Cast {
		if isSeen(cast.Url) {
			actors = append(actors, cast.Url)
			if cast.Url == destination {
				fmt.Println("DONE --> ", cast.Url)
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
