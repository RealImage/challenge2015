package main

import (
	"fmt"
	"os"
)

const moviebuff = "http://data.moviebuff.com/"

var (
	actorList map[string][]string
	seen      map[string]bool
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
	actorList = make(map[string][]string)
	buildTree(os.Args[1], os.Args[2])
}

func buildTree(argument, destination string) {
	degrees++
	actorList[argument] = []string{}
	buildActors(argument, argument, destination)
	fmt.Println(actorList[argument])
	for _, actor := range actorList[argument]{
		buildTree(actor, destination)
	}
}

func buildActors(argument, parent, destination string) {
	url := moviebuff + argument
	json, err := getData(url)
	defer ErrHandle(err)

	for _, movie := range json.Movies {
		if isSeen(movie.Url) {
			buildActors(movie.Url, argument, destination)
		}
	}

	for _, cast := range json.Cast {
		if isSeen(cast.Url) {
			actorList[parent] = append(actorList[parent], cast.Url)
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
