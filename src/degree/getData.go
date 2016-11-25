package main

import (
	"encoding/json"
	//	"fmt"
	"io/ioutil"
	"net/http"
	//	"os"
	//	"runtime"

	"time"

	"github.com/beefsack/go-rate"
)

// Set the rate of operations/sec
var rl = rate.New(100, time.Second)

// Get the movie details using the movie url
func getMovieData(url string) (movie, error) {
	for {
		if ok, _ := rl.Try(); ok {
			// Get the movies data
			re, err := http.Get("http://data.moviebuff.com/" + url)
			if err != nil {
				//				fmt.Println("invalid url:", url, "Error:", err.Error())
				return movie{}, err
				//// Retry
				//time.Sleep(10 * time.Millisecond)
				//continue
			}
			// Read the json data
			data, err := ioutil.ReadAll(re.Body)
			if err != nil {
				//				fmt.Println("invalid url:", url, "Error:", err.Error())
				return movie{}, err
				//// Retry
				//time.Sleep(10 * time.Millisecond)
				//continue
			}
			var m movie
			// Unmarshal the json data to movie struct
			err = json.Unmarshal(data, &m)
			if err != nil {
				//				fmt.Println("invalid url:", url, "Error:", err.Error())
				return movie{}, err
				//// Retry
				//time.Sleep(10 * time.Millisecond)
				//continue
			}
			return m, err
		} else {
			time.Sleep(10 * time.Millisecond)
		}
	}
}

// Get person details using his/her url
func getPersonData(url string) (person, error) {
	for {
		if ok, _ := rl.Try(); ok {
			// Get the person data
			re, err := http.Get("http://data.moviebuff.com/" + url)
			if err != nil {
				//				fmt.Println("invalid url:", url, "Error:", err.Error())
				return person{}, err
				//// Retry
				//time.Sleep(10 * time.Millisecond)
				//continue
			}
			// Read the json data
			data, err := ioutil.ReadAll(re.Body)
			if err != nil {
				//				fmt.Println("invalid url:", url, "Error:", err.Error())
				return person{}, err
				//// Retry
				//time.Sleep(10 * time.Millisecond)
				//continue
			}
			var p person
			// Unmarshal the json data to person struct
			err = json.Unmarshal(data, &p)
			if err != nil {
				//				fmt.Println("invalid url:", url, "Error:", err.Error())
				return person{}, err
				//// Retry
				//time.Sleep(10 * time.Millisecond)
				//continue
			}
			return p, err
		} else {
			time.Sleep(1 * time.Millisecond)
		}
	}
}
