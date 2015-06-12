package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var URL = "https://data.moviebuff.com/"

var client *http.Client

func init() {
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, DisableKeepAlives: false, MaxIdleConnsPerHost: 20}
	client = &http.Client{Transport: tr}
}

func GetPersonDetails(id string) *Person {
	FileLogger.Println("GetPersonDetails::", id)
	url := fmt.Sprint(URL, id)

	resp, err := client.Get(url)
	if err != nil {
		// TODO
		FileLogger.Println("got some error 1")
		FileLogger.Println(err)
		somePerson := Person{Url: id, Name: id}
		return &somePerson
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//TODO
		FileLogger.Println("got some error 2")
		FileLogger.Println(err)
		somePerson := Person{Url: id, Name: id}
		return &somePerson
	}

	personData := Person{Movies: make([]PersonMovies, 0)}
	err = json.Unmarshal(respData, &personData)
	if err != nil {
		//TODO
		FileLogger.Println(err)
		FileLogger.Println("got some error 3")
		personData.Name = id
	}
	FileLogger.Println("returning from GetPersonDetails")
	return &personData
}

func GetMovieDetails(id string) *Movie {
	FileLogger.Println("GetMovieDetails::", id)
	url := fmt.Sprint(URL, id)

	resp, err := client.Get(url)
	if err != nil {
		// TODO
		FileLogger.Println(err)
		FileLogger.Println("got some movie error 1")
		someMovie := Movie{Url: id, Name: id}
		return &someMovie
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//TODO
		FileLogger.Println(err)
		FileLogger.Println("got some movie error 2")
		someMovie := Movie{Url: id, Name: id}
		return &someMovie
	}

	movieData := Movie{Cast: make([]MovieCast, 0), Crew: make([]MovieCrew, 0)}
	err = json.Unmarshal(respData, &movieData)
	if err != nil {
		//TODO
		FileLogger.Println(err)
		FileLogger.Println("got some movie error 3")
		movieData.Name = id
	}
	FileLogger.Println("returning from GetMovieDetails")
	return &movieData
}
