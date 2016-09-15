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
	retList := make(map[string][]string)

	degrees++
	retList = loopMovies(os.Args[1], os.Args[1], os.Args[2])

	if len(retList) != 0 {
		degrees++
		fmt.Println(len(retList))
		for k := range retList {
			for _, v := range retList[k] {
				fmt.Println(v)
				loopMovies(v, v, os.Args[2])
			}
		}
	}
}

func loopMovies(argument, parent, destination string) map[string][]string {
	retList := make(map[string][]string)
	url := moviebuff + argument
	json, err := getData(url)
	defer ErrHandle(err)
	for _, movie := range json.Movies {
		if notSeen(movie.Url) {
			retList[argument] = loopActors(movie.Url, argument,
				destination, retList[argument])
		}
	}
	return retList
}

func loopActors(argument, parent, destination string, retList []string) []string {
	url := moviebuff + argument
	json, err := getData(url)
	defer ErrHandle(err)
	for _, cast := range json.Cast {
		if notSeen(cast.Url) {
			retList = append(retList, cast.Url)
			if cast.Url == destination {
				fmt.Println("DONE --> ", cast.Url, degrees)
				os.Exit(1)
			}
		}
	}
	return retList

}

func notSeen(in string) bool {
	if !seen[in] {
		seen[in] = true
		return true
	} else {
		return false
	}
}
