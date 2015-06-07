package main

import (
	"fmt"
	"strings"
)

var dictSource map[string]string
var dictDestination map[string]string

var SetSource map[string]bool
var SetDestination map[string]bool
var SetSrcMovie map[string]bool
var SetDestMovie map[string]bool

var depthNamesSrc map[int]string
var depthNamesDest map[int]string

var depth int
var newNames string

func mapInitialize() {
	dictSource = make(map[string]string)
	dictDestination = make(map[string]string)

	SetSource = make(map[string]bool)
	SetDestination = make(map[string]bool)

	SetSrcMovie = make(map[string]bool)
	SetDestMovie = make(map[string]bool)

	depthNamesSrc = make(map[int]string)
	depthNamesDest = make(map[int]string)

	depth = 0
}

func GetSeparation(fromId string, toId string) {
	FileLogger.Println("inside GetSeparation")
	mapInitialize()
	depthNamesSrc[0] = fromId
	depthNamesDest[0] = toId

	sourcePath := fmt.Sprint(fromId, "~")
	destinationPath := fmt.Sprint(toId, "~")
	dictSource[fromId] = "0"
	dictDestination[toId] = "0"

	allSrcNames := ""
	allDestNames := ""
	srcKey := ""
	srcValue := ""
	destValue := ""

Loop:
	for {

		allSrcNames = depthNamesSrc[depth]
		allDestNames = depthNamesDest[depth]
		depth++
		sNames := strings.Split(allSrcNames, "~")
		for _, v := range sNames {
			_, done := SetSource[v]
			if !done {
				person := GetPersonDetails(v)
				SetSource[person.Url] = true
				res := process(person, sourcePath, "source")
				if res != nil {
					srcKey = res.key
					srcValue = res.srcValue
					destValue = res.destValue
					break Loop
				}
			}
		}

		dNames := strings.Split(allDestNames, "~")
		for _, v := range dNames {
			_, done := SetDestination[v]
			if !done {
				person := GetPersonDetails(v)
				SetDestination[person.Url] = true
				res := process(person, destinationPath, "destination")
				if res != nil {
					srcKey = res.key
					srcValue = res.srcValue
					destValue = res.destValue
					break Loop
				}
			}
		}
	}
	FileLogger.Println("srcKey::", srcKey)
	displayResult(srcValue, destValue)
	FileLogger.Println("returning from GetSeparation")
}

func getActorsFromMovies(movies []PersonMovies, tag string, path string) {
	FileLogger.Println("inside getActorsFromMovies")
	newNames = ""
	for _, movie := range movies {
		if tag == "source" {
			_, found := SetSrcMovie[movie.Url]
			if found {
				FileLogger.Println("this movie already processed for src")
				return
			} else {
				SetSrcMovie[movie.Url] = true
			}
		} else {
			_, found := SetDestMovie[movie.Url]
			if found {
				FileLogger.Println("this movie already processed for dest")
				return
			} else {
				SetDestMovie[movie.Url] = true
			}
		}
		movieDetailsptr := GetMovieDetails(movie.Url)
		getNames(movieDetailsptr, tag, path)
	}
	if tag == "source" {
		depthNamesSrc[depth] = newNames
	} else {
		depthNamesDest[depth] = newNames
	}
}

func getNames(movie *Movie, tag string, path string) {
	FileLogger.Println("inside getNames")
	for _, cast := range movie.Cast {
		url := cast.Url
		if tag == "source" {
			dictSource[url] = fmt.Sprint(path, movie.Url, "~", url, "~")
		} else {
			dictDestination[url] = fmt.Sprint(path, movie.Url, "~", url, "~")
		}
		newNames = fmt.Sprint(newNames, url, "~")
	}
	for _, crew := range movie.Crew {
		url := crew.Url
		if tag == "source" {
			dictSource[url] = fmt.Sprint(path, movie.Url, "~", url, "~")
		} else {
			dictDestination[url] = fmt.Sprint(path, movie.Url, "~", url, "~")
		}
		newNames = fmt.Sprint(newNames, url, "~")
	}
}

func verify() (*match, bool) {
	FileLogger.Println("inside verify")
	found := false
	res := match{}
	for key := range dictDestination {
		sourceValue, exists := dictSource[key]
		if exists {
			FileLogger.Println("oh GOD..found the break point.")
			FileLogger.Println("key::", key)
			FileLogger.Println("sourceValue::", sourceValue)
			FileLogger.Println("dest value::", dictDestination[key])
			res.key = key
			res.srcValue = sourceValue
			res.destValue = dictDestination[key]
			found = true
			break
		}
	}
	return &res, found
}

func process(person *Person, path string, tag string) *match {
	FileLogger.Println("inside process")
	getActorsFromMovies(person.Movies, tag, path)
	FileLogger.Println("added new names")
	res, found := verify()
	if found {
		return res
	}
	return nil
}

func displayResult(srcValue string, destValue string) {
	FileLogger.Println("inside displayResult")
	srcIds := strings.Split(srcValue, "~")
	destIds := strings.Split(destValue, "~")
	movies := 0
	people := 0
	totalLength := len(srcIds) + len(destIds)
	a := make([]string, totalLength)
	b := make([]string, totalLength)
	for i, v := range srcIds {
		v = strings.TrimSpace(v)
		if v != "" {
			if i%2 == 0 {
				a[people] = v
				people++
			} else {
				b[movies] = v
				movies++
			}
		}
	}
	j := len(destIds) - 3
	FileLogger.Println("len::", len(destIds))
	FileLogger.Println("j::", j)
	i := j
	for ; i >= 0; i-- {
		v := destIds[i]
		v = strings.TrimSpace(v)
		if v != "" {
			if i%2 == 0 {
				a[people] = destIds[i]
				people++
			} else {
				b[movies] = destIds[i]
				movies++
			}
		}
	}
	FileLogger.Println("a::", a)
	FileLogger.Println("b::", b)
	fmt.Println("Degrees of Separation:", movies)
	j = 0
	p := GetPersonDetails(a[0])
	found := false
	for i, v := range b {
		v = strings.TrimSpace(v)
		if v != "" {
			x := GetMovieDetails(v)
			fmt.Printf("%d. Movie: %s\n", i+1, x.Name)

			found = false
			for _, n := range x.Cast {
				if n.Url == a[j] {
					fmt.Println(n.Role, ":", p.Name)
					found = true
					break
				}
			}
			if !found {
				for _, m := range x.Crew {
					if m.Url == a[j] {
						fmt.Println(m.Role, ":", p.Name)
						break
					}
				}
			}
			j++
			p = GetPersonDetails(a[j])
			found = false
			for _, n := range x.Cast {
				if n.Url == a[j] {
					fmt.Println(n.Role, ":", p.Name)
					found = true
					break
				}
			}
			if !found {
				for _, m := range x.Crew {
					if m.Url == a[j] {
						fmt.Println(m.Role, ":", p.Name)
						break
					}
				}
			}

		}
	}
	FileLogger.Println("returning from displayResult")
}
