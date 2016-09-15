package main

import (
	"fmt"
	"os"
)

const moviebuff = "http://data.moviebuff.com/"

var (
	seen      map[string]bool
	degrees   int
	destination string
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
	var q queue

	degrees++
	retList[os.Args[1]] = loopMovies(os.Args[1], os.Args[1], os.Args[2])

	for k := range retList{
		q.enqueue(k)
	}
	for len(q.value) != 0 {
		degrees++
		fmt.Println("Looking into level ",  degrees)
		fmt.Println(q.value)
		for _, k := range q.value{
		q.dequeue()
			for _, v := range retList[k] {
				fmt.Println(v)
				retList[v] = loopMovies(v, v, os.Args[2])
				q.enqueue(v)
			}
		}
	}
}

func loopMovies(argument, parent, destination string) []string {
	var retList []string
	url := moviebuff + argument
	json, err := getData(url)
	defer ErrHandle(err)
	for _, movie := range json.Movies {
		if notSeen(movie.Url) {
			retList = loopActors(movie.Url, argument,
				destination, retList)
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
	for _, crew := range json.Crew {
		if notSeen(crew.Url) {
			retList = append(retList, crew.Url)
			if crew.Url == destination {
				fmt.Println("DONE --> ", crew.Url, degrees)
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
