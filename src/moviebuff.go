package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var URL = "https://data.moviebuff.com/"

func GetPersonDetails(id string) *Person {
	FileLogger.Println("GetPersonDetails::", id)
	url := fmt.Sprint(URL, id)
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		// TODO
		FileLogger.Println(err)
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//TODO
		FileLogger.Println(err)
	}

	personData := Person{Movies: make([]PersonMovies, 0)}
	err = json.Unmarshal(respData, &personData)
	if err != nil {
		//TODO
		FileLogger.Println(err)
		personData.Name = id
	}
	FileLogger.Println("returning from GetPersonDetails")
	return &personData
}

func GetMovieDetails(id string) *Movie {
	FileLogger.Println("GetMovieDetails::", id)
	url := fmt.Sprint(URL, id)
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		// TODO
		FileLogger.Println(err)
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//TODO
		FileLogger.Println(err)
	}

	movieData := Movie{Cast: make([]MovieCast, 0), Crew: make([]MovieCrew, 0)}
	err = json.Unmarshal(respData, &movieData)
	if err != nil {
		//TODO
		FileLogger.Println(err)
		movieData.Name = id
	}
	FileLogger.Println("returning from GetMovieDetails")
	return &movieData
}
