package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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

/* For back trace */
type Path struct {
	name  string
	movie Movie
}

var trace map[string]Path
var processedMovie map[string]bool
var processedPerson map[string]bool
var degree int

func formUrl(url string) string {
	url = homeurl + url
	return url
}

func getRoleFromMovie(movie string, person string) CCtype {
	url := formUrl(movie)
	var cc CastCrew
	var ccrole CCtype

	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return ccrole
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
	for _, ccrole = range cc.Cast {
		if ccrole.Url == person {
			goto tag
		}
	}
	for _, ccrole = range cc.Crew {
		if ccrole.Url == person {
			goto tag
		}
	}
tag:
	return ccrole
}

func backtrace(personA string, personB string) {
	var prevname string
	path := make([]Path, 0)
	path = append(path, trace[personB])
	br := 0

	for path[len(path) - 1].name != personA && br < degree-1 {
		path = append(path, trace[path[len(path) - 1].name])
		br++
	}

	/*Print Path */
	for i, v := range path {
		fmt.Println("Movie:", v.movie.Url)
		if i == 0 {
			role := getRoleFromMovie(path[0].movie.Url, personB)
			fmt.Println(role.Role, ":", personB)	
		} else {
			role := getRoleFromMovie(v.movie.Url, prevname) 	
			fmt.Println(role.Role, ":", prevname)
		}
			
		fmt.Println(v.name, ":", v.movie.Role, "\n")
		prevname = v.name
	}
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

func getCCfromMovies(movie_list []Movie, person string, personB string) ([]CCtype, bool) {
	cclist := make([]CCtype, 0)
	var path Path

	for _, val := range movie_list {
		//fmt.Println("CastCrew", i, val)
		cc, err := getCastCrew(val.Url)
		if err != nil {
			continue
		}
		cclist = append(cclist, cc.Cast...)
		cclist = append(cclist, cc.Crew...)
		for _, val2 := range cc.Cast {
			path.name = person
			path.movie = val
			trace[val2.Url] = path
			if val2.Url == personB {
				return cclist, true
			}
		}
		for _, val2 := range cc.Crew {
			path.name = person
			path.movie = val
			trace[val2.Url] = path
			if val2.Url == personB {
    			return cclist, true
			}
		}
	}
	return cclist, false
}

func getCCfromCC(person_list []CCtype, personB string) ([]CCtype, bool) {
	cclist := make([]CCtype, 0)
	var tmp Path

	for _, person := range person_list {
		//fmt.Println("Person", i, person)
		mval, ok := processedPerson[person.Url]
		if ok == false || mval != true {
			movie_list, err := getPersonMovies(person.Url)
			if err != nil {
				continue
			}
			for _, movie := range movie_list {
				mval, ok := processedMovie[movie.Url]
				if ok == false || mval != true {
					cc, _ := getCastCrew(movie.Url)
					for _, val := range cc.Cast {
						tmp.name = person.Url
						tmp.movie = movie
						trace[val.Url] = tmp
						if val.Url == personB {
							//fmt.Println("\n", val, personB, "\n")
							return cclist, true
						}
					}
					for _, val := range cc.Crew {
						tmp.name = person.Url
						tmp.movie = movie
						trace[val.Url] = tmp
						if val.Url == personB {
							//fmt.Println("\n", val, personB, "\n")
     						return cclist, true
 						}
					}
					cclist = append(cclist, cc.Cast...)
					cclist = append(cclist, cc.Crew...)
					processedMovie[movie.Url] = true
				}
			}
		}
		processedPerson[person.Url] = true
	}
	return cclist, false
}

func processPerson(personA string, personB string) {
	degree = 1
	movie_list, _ := getPersonMovies(personA)
	castcrew_list, rc := getCCfromMovies(movie_list, personA, personB)
	if rc == true {
		fmt.Println("\nDegree of Seperation -", degree, "\n")
		return 
	}
	for degree <= 6 {
		degree++
		castcrew_list, rc = getCCfromCC(castcrew_list, personB)
		if rc == true {
			fmt.Println("\nDegree of Seperation -", degree, "\n")
			return
		}
	}
	fmt.Println("Exceed the degree")
}

func main() {
	args := os.Args
	trace = make(map[string]Path)
	processedMovie = make(map[string]bool)
	processedPerson = make(map[string]bool)

	if len(args) < 3 {
		fmt.Println("Invalid Arguments")
		os.Exit(1)
	}
	personA := args[1]
	personB := args[2]
	processPerson(personA, personB)
	backtrace(personA, personB)
}
