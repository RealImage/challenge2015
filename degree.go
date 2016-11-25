package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	homeurl           = "http://data.moviebuff.com/"
	routines_per_core = 10
	url_retry         = 10
	debug             = false
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
var m sync.RWMutex

var processedMovie map[string]bool
var processedPerson map[string]bool
var degree int

func Log(v ...interface{}) {
	if debug == true {
		fmt.Println(v)
	}
}

func formUrl(url string) string {
	url = homeurl + url
	return url
}

func getUrlData(url string) ([]byte, error) {
	var urldata []byte
	var err error
	it := 0
	url = formUrl(url)

	for it < url_retry {
		res, err := http.Get(url)
		if err != nil {
			Log("Retrying Url:", url)
			time.Sleep(10 * time.Millisecond)
			it++
			continue
		} else {
			defer res.Body.Close()
			urldata, err = ioutil.ReadAll(res.Body)
			break
		}
	}
	return urldata, err
}

func getRoleFromMovie(movie string, person string) CCtype {
	var cc CastCrew
	var ccrole CCtype

	urldata, err := getUrlData(movie)

	err = json.Unmarshal(urldata, &cc)
	if err != nil {
		Log("JSON Decode:", err)
		err = nil
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

	for path[len(path)-1].name != personA && br < degree-1 {
		path = append(path, trace[path[len(path)-1].name])
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
	var person Person

	urldata, err := getUrlData(name)

	err = json.Unmarshal(urldata, &person)
	if err != nil {
		//Log("JSON Decode:", err)
		err = nil //just ignore
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
	var cc CastCrew

	urldata, err := getUrlData(name)

	err = json.Unmarshal(urldata, &cc)
	if err != nil {
		//Log("JSON Decode:", err)
		err = nil 
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

func getCCfromCC2(person_list []CCtype, personB string, ch_list chan []CCtype) {
	cclist := make([]CCtype, 0)
	var tmp Path

	for _, person := range person_list {
		//Log(len(person_list), "- Person", i, person)
		movie_list, err := getPersonMovies(person.Url)
		if err != nil {
			continue
		}
		for _, movie := range movie_list {
			cc, _ := getCastCrew(movie.Url)
			for _, val := range cc.Cast {
				tmp.name = person.Url
				tmp.movie = movie
				m.Lock()
				trace[val.Url] = tmp
				m.Unlock()
			}
			for _, val := range cc.Crew {
				tmp.name = person.Url
				tmp.movie = movie
				m.Lock()
				trace[val.Url] = tmp
				m.Unlock()
			}

			cclist = append(cclist, cc.Cast...)
			cclist = append(cclist, cc.Crew...)
		}
	}
	ch_list <- cclist
}

func getCCfromCC(person_list []CCtype, personB string) ([]CCtype, bool) {
	var list_count, cycles, rem_list int

	routines := runtime.NumCPU() * routines_per_core
	list_count = len(person_list)

	if list_count > routines {
		cycles = list_count / routines
		rem_list = list_count % routines
	} else {
		routines = 0
		cycles = 0
		rem_list = list_count
	}
	Log("No Routines-", routines, "Total-", list_count, "cycles-", cycles)

	ch_list := make(chan []CCtype)
	cclist := make([]CCtype, 0)
	found := false

	defer close(ch_list)
	doneCount := 0
	for it := 0; it < routines; it++ {
		Log("Routine", it, "started")
		go getCCfromCC2(person_list[(it*cycles):cycles+(it*cycles)], personB, ch_list)
	}
	if rem_list > 0 {
		Log("Routine started- rem")
		go getCCfromCC2(person_list[cycles*routines:], personB, ch_list)
		routines++
	}
	for doneCount != routines {
		select {
		case cc_list := <-ch_list:
			Log("Routine", doneCount, "finished")
			cclist = append(cclist, cc_list...)
			doneCount++
		}
	}

	for _, val := range cclist {
		if val.Url == personB {
			found = true
			break
		}
	}
	return cclist, found
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
	runtime.GOMAXPROCS(runtime.NumCPU())
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

/*func getCCfromCC(person_list []CCtype, personB string) ([]CCtype, bool) {
	cycles := cycles_per_routine
	list_count := len(person_list)
	routine_count := list_count / cycles
	rem_list := list_count % cycles
	Log("Number of  elements in list", list_count)
	Log("Number of go routines", routine_count)


	ch_list := make(chan []CCtype)
	cclist := make([]CCtype, 0)
	found := false

	for it := 0; it < routine_count; it++ {
		Log("Routine", it, "started")
		go getCCfromCC2(person_list[(it*cycles):cycles+(it*cycles)], personB, ch_list)
	}
	for it := 0; it < routine_count; it++ {
		cc_list := <-ch_list
		fmt.Println("Routine", it, "finished")
		cclist = append(cclist, cc_list...)
	}
	if rem_list > 0 {
		fmt.Println("Routine started- rem")
		go getCCfromCC2(person_list[routine_count*cycles:], personB, ch_list)
		cc_list := <-ch_list
		cclist = append(cclist, cc_list...)
		fmt.Println("Routine finished: rem")
	}
	for _, val := range cclist {
		if val.Url == personB {
			found = true
			break
		}
	}
	return cclist, found
}*/
