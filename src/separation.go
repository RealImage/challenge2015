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

	dictSource[fromId] = fmt.Sprint(fromId, "~")
	dictDestination[toId] = fmt.Sprint(toId, "~")

	allSrcNames := ""
	allDestNames := ""
	srcKey := ""
	srcValue := ""
	destValue := ""

Loop:
	for {

		allSrcNames = depthNamesSrc[depth]
		allDestNames = depthNamesDest[depth]
		FileLogger.Println("depth::", depth)
		depth++
		sNames := strings.Split(allSrcNames, "~")
		chSrc := make(chan match)
		chDest := make(chan match)
		for _, v := range sNames {
			go callProcess(v, "source", chSrc)
		}

		dNames := strings.Split(allDestNames, "~")
		for _, v := range dNames {
			go callProcess(v, "destination", chDest)
		}

		for range sNames {
			res := <-chSrc
			if res.key != "" {
				srcKey = res.key
				srcValue = res.srcValue
				destValue = res.destValue
				break Loop
			}
		}
		for range dNames {
			res := <-chDest
			if res.key != "" {
				srcKey = res.key
				srcValue = res.srcValue
				destValue = res.destValue
				break Loop
			}
		}
		FileLogger.Println("calling verify before going more depth")
		mat, found := verify()
		if found {
			srcKey = mat.key
			srcValue = mat.srcValue
			destValue = mat.destValue
			break Loop
		}
	}
	FileLogger.Println("srcKey::", srcKey)
	displayResult(srcValue, destValue)
	FileLogger.Println("returning from GetSeparation")
}

func callProcess(aName string, tag string, someCh chan match) {
	FileLogger.Println("inside callProcess")
	var res match
	if tag == "source" {
		_, done := SetSource[aName]
		if !done {
			person := GetPersonDetails(aName)
			SetSource[person.Url] = true
			sourcePath := dictSource[person.Url]
			res, _ = process(person, sourcePath, tag)
		}
	} else {
		_, done := SetDestination[aName]
		if !done {
			person := GetPersonDetails(aName)
			SetDestination[person.Url] = true
			destinationPath := dictDestination[person.Url]
			res, _ = process(person, destinationPath, tag)

		}
	}
	someCh <- res
	FileLogger.Println("returning from callProcess")
}

func getActorsFromMovies(movies []PersonMovies, tag string, path string) {
	FileLogger.Println("inside getActorsFromMovies")
	ch := make(chan string)
	for _, movie := range movies {

		go movie1(movie.Url, tag, path, ch)
	}
	for range movies {
		aName := <-ch
		if tag == "source" {
			FileLogger.Println("before assigning src")
			//FileLogger.Println(depthNamesSrc[depth])
			depthNamesSrc[depth] = fmt.Sprint(depthNamesSrc[depth], aName)
		} else {
			FileLogger.Println("before assigning dest")
			//FileLogger.Println(depthNamesDest[depth])
			depthNamesDest[depth] = fmt.Sprint(depthNamesDest[depth], aName)
		}

	}

}

func movie1(movieUrl string, tag string, path string, ch chan string) {
	FileLogger.Println("inside movie1")
	processed := false
	if tag == "source" {
		_, found := SetSrcMovie[movieUrl]
		if found {
			FileLogger.Println("this movie already processed for src")
			processed = true
		} else {
			SetSrcMovie[movieUrl] = true
		}
	} else {
		_, found := SetDestMovie[movieUrl]
		if found {
			FileLogger.Println("this movie already processed for dest")
			processed = true
		} else {
			SetDestMovie[movieUrl] = true
		}
	}
	if !processed {
		movieDetailsptr := GetMovieDetails(movieUrl)
		names := getNames(movieDetailsptr, tag, path)
		ch <- names
	} else {
		ch <- ""
	}
}

func getNames(movie *Movie, tag string, path string) string {
	FileLogger.Println("inside getNames")
	names := ""
	for _, cast := range movie.Cast {
		url := cast.Url
		if tag == "source" {
			dictSource[url] = fmt.Sprint(path, movie.Url, "~", url, "~")
		} else {
			dictDestination[url] = fmt.Sprint(path, movie.Url, "~", url, "~")
		}
		names = fmt.Sprint(names, url, "~")
	}
	for _, crew := range movie.Crew {
		url := crew.Url
		if tag == "source" {
			dictSource[url] = fmt.Sprint(path, movie.Url, "~", url, "~")
		} else {
			dictDestination[url] = fmt.Sprint(path, movie.Url, "~", url, "~")
		}
		names = fmt.Sprint(names, url, "~")
	}
	return names
}

func verify() (match, bool) {
	FileLogger.Println("inside verify")
	found := false
	res := match{}
	for key := range dictDestination {
		sourceValue, exists := dictSource[key]
		if exists {
			FileLogger.Println("oh GOD..found the link.")
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
	return res, found
}

func process(person *Person, path string, tag string) (match, bool) {
	FileLogger.Println("inside process")
	getActorsFromMovies(person.Movies, tag, path)
	FileLogger.Println("added new names")
	res, found := verify()
	return res, found
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
	FileLogger.Println("names::", a)
	FileLogger.Println("movies::", b)
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
