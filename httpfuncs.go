package main

import (
	//"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	//"sort"
)

const moviebuff = "http://data.moviebuff.com/"

/*
func filterQueue(in chan string, out chan string) {
	var seen = make(map[string]bool)
	for val := range in {
		if !seen[val] {
			seen[val] = true
			out <- val
		}
	}
}
*/

func ErrHandle(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

/*
func contains(in string, list []string) bool {
	sort.Strings(list)
	i := sort.SearchStrings(list, in)
	return (i < len(list) && list[i] == in)
}

func main() {

	url := moviebuff + os.Args[1]

	queue := make(chan string)
	filteredQueue := make(chan string)
	go func() {
		queue <- url
	}()

	go filterQueue(queue, filteredQueue)

	for newurl := range filteredQueue {
		fmt.Println(newurl)
		enqueue(newurl, queue, os.Args[2])
	}

}

func enqueue(url string, queue chan string, final string) {

	var Head header

	resp, err := http.Get(url)
	ErrHandle(err)
	body, err := ioutil.ReadAll(resp.Body)
	ErrHandle(err)
	err = json.Unmarshal(body, &Head)
	switch Head.Typ {

	case "Movie":
		persons, err := getCast(body)
		ErrHandle(err)

		for _, person := range persons {
			if person.Url == final{
				fmt.Println(person.Url, url)
				os.Exit(1)
			}
		}

		for _, person := range persons {
			go func() {
				queue <- moviebuff + person.Url
			}()
		}
	case "Person":
		movies, err := getMovies(body)
		ErrHandle(err)
		//		for _, movie := range movies {
		//			fmt.Println(movie.Url, Head.Url)
		//		})
		//
		for _, movie := range movies {
			go func() {
				queue <- moviebuff + movie.Url
			}()

		}

	}
}
*/

func main(){
	url := moviebuff + os.Args[1]
	resp, err := http.Get(url)
	defer ErrHandle(err)
	body, err := ioutil.ReadAll(resp.Body)

	json, err := getData(body)
	fmt.Println(json.Name, json.Url)

}
