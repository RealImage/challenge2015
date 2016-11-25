package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"time"

	"github.com/beefsack/go-rate"
)

// Set the rate of operations/sec
var rl *rate.RateLimiter

// Get the movie details using the movie url
func getMovieData(url string) (movie, error) {
	if rl == nil {
		rl = rate.New(conf.Rate, time.Second)
	}
	for {
		if ok, _ := rl.Try(); ok {
			// Get the movies data
			re, err := http.Get(conf.Url + url)
			if err != nil {
				return movie{}, err
			}
			// Read the json data
			data, err := ioutil.ReadAll(re.Body)
			if err != nil {
				return movie{}, err
			}
			var m movie
			// Unmarshal the json data to movie struct
			err = json.Unmarshal(data, &m)
			if err != nil {
				return movie{}, err
			}
			return m, err
		} else {
			time.Sleep(10 * time.Millisecond)
		}
	}
}

// Get person details using his/her url
func getPersonData(url string) (person, error) {
	if rl == nil {
		rl = rate.New(conf.Rate, time.Second)
	}
	for {
		if ok, _ := rl.Try(); ok {
			// Get the person data
			re, err := http.Get(conf.Url + url)
			if err != nil {
				return person{}, err
			}
			// Read the json data
			data, err := ioutil.ReadAll(re.Body)
			if err != nil {
				return person{}, err
			}
			var p person
			// Unmarshal the json data to person struct
			err = json.Unmarshal(data, &p)
			if err != nil {
				return person{}, err
			}
			return p, err
		} else {
			time.Sleep(1 * time.Millisecond)
		}
	}
}
