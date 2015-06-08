package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	//"time"
)

const (
	homeurl   = "http://data.moviebuff.com/"
	url_retry = 2
)

/* For Read purpose same type struct with diff name */
type Movie struct {
	Name string
	Url  string
	Role string
}

type CCtype struct {
	Url  string
	Name string
	Role string
}

type CastCrew struct {
	Url  string   `json:"url"`
	Type string   `json:"type"`
	Name string   `json:"name"`
	Cast []CCtype `json:"cast"`
	Crew []CCtype `json:"crew"`
}

type Person struct {
	Url    string  `json:"url"`
	Type   string  `json:"type"`
	Name   string  `json:"name"`
	Movies []Movie `json:"movies"`
}

func formUrl(url string) string {
	url = homeurl + url
	return url
}


func getCastCrew(name string) (CastCrew, error) {
	url := formUrl(name)
	var cc CastCrew

	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return cc, err
	}
	defer res.Body.Close()

	urldata, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("io read err", err)
	}

	err = json.Unmarshal(urldata, &cc)
	if err != nil {
		//fmt.Println("json decode error", err)
	}

	return cc, err
}

func getCastCrewList(movie_list []Movie) []CCtype {
	cclist := make([]CCtype, 0)

	for _, val := range movie_list {
		//fmt.Println ("CastCrew", i, val)
		cc, err := getCastCrew(val.Url)
		if err != nil {
			continue
		}
		cclist = append(cclist, cc.Cast...)
		cclist = append(cclist, cc.Crew...)
	}
	return cclist
}

func getPersonMovies(name string) ([]Movie, error) {
	url := formUrl(name)
	var person Person

	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return person.Movies, err
	}
	defer res.Body.Close()

	urldata, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("io read err", err)
	}

	err = json.Unmarshal(urldata, &person)
	if err != nil {
		//fmt.Println("json decode error", err)
	}
	return person.Movies, err
}

func getPersonMovieList(person_list []CCtype) []Movie {
	movie_list := make([]Movie, 0)

	for _, val := range person_list {
		//fmt.Println("Person", i, val)
		movies, err := getPersonMovies(val.Url)
		if err != nil {
			continue
		}
		movie_list = append(movie_list, movies...)
	}
	return movie_list

}

func processPerson(personA string, personB string) {
	degree := 1
	movie_list, _ := getPersonMovies(personA)
	castcrew_list := getCastCrewList(movie_list)
	for degree <= 6 {
		for _, cc := range castcrew_list {
			//fmt.Println(i, cc)
			if cc.Url == personB {
				fmt.Println("Degree of Seperation Between\n\t", personA, "&", personB, "-", degree)
				return
			}
		}
		fmt.Println("Not in degree ", degree)
		movie_list = getPersonMovieList(castcrew_list)
		castcrew_list = getCastCrewList(movie_list)
		degree++
	}
	fmt.Println("Exceed the degree")
}

func main() {
	args := os.Args

	if len(args) < 3 {
		fmt.Println("Invalid Arguments")
		os.Exit(1)
	}
	personA := args[1]
	personB := args[2]
	processPerson(personA, personB)
}
