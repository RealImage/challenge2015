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
	//var queue []string
	if len(actorList) == 0 {
		degrees++
		loopMovies(os.Args[1], os.Args[1], os.Args[2])
	}
	if len(actorList) != 0 {
		for true {
			degrees++
			for _, v := range actorList[os.Args[1]] {
				fmt.Println(v)
				loopMovies(v,v, os.Args[2])
			}
		}
	}
}

func buildTree(argument, destination string) {
	//var retList map[string][]string
	//retList[argument] = append(retList[argument]),loopMovies(argument, argument, destination))
}

func loopMovies(argument, parent, destination string) {
	//var retList []string
	url := moviebuff + argument
	json, err := getData(url)
	defer ErrHandle(err)
	for _, movie := range json.Movies {
		if notSeen(movie.Url) {
			loopActors(movie.Url, argument, destination)
		}
	}
}

func loopActors(argument, parent, destination string){
	url := moviebuff + argument
	json, err := getData(url)
	defer ErrHandle(err)
	for _, cast := range json.Cast {
		if notSeen(cast.Url) {
			actorList[parent] = append(actorList[parent], cast.Url)
			if cast.Url == destination {
				fmt.Println("DONE --> ", cast.Url, degrees)
				os.Exit(1)
			}
		}
	}

}

func notSeen(in string) bool {
	if !seen[in] {
		seen[in] = true
		return true
	} else {
		return false
	}
}
