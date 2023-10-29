
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/RealImage/Challenge/models"
	"github.com/RealImage/Challenge/utilities"
)

var (
	Log *log.Logger
)

func init() {
	file, err := os.Create("info.log")
	if err != nil {
		panic(err)
	}

	Log = log.New(file, "", log.LstdFlags|log.Lshortfile)
}

func validateInputs() bool {
	res := make(chan models.Response)
	throttle := time.Tick(time.Second / 800)

	go getdata(os.Args[1], res, throttle)

	r := <-res
	if r.Err != nil {
		return false
	}

	go getdata(os.Args[2], res, throttle)

	r = <-res
	if r.Err != nil {
		return false
	}

	return true
}

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Usage: degrees <actor-name> <actor-name>")
	}

	ok := validateInputs()
	if !ok {
		log.Fatal("Invalid inputs")
	}

	out := make(chan models.Out)
	q := utilities.NewQueue()
	q.Enqueue(utilities.WorkRequest{Actor: os.Args[1]})

	go getrelationship(q, out)

	for {
		o := <-out
		if o.Err != nil {
			if _, ok := o.Err.(models.Skippable); ok {
				continue
			}
			panic(o.Err)
		}

		if o.Relationship == nil {
			continue
		}

		if o.Relationship.Cast2.URL == os.Args[2] {
			var result []models.Relationship
			result = o.Relationship.Path
			result = append(result, *o.Relationship)

			fmt.Println("Degrees of seperation:", len(result))
			fmt.Println()

			for i, r := range result {
				fmt.Printf("%d. Movie: %s\n", i+1, r.Movie)
				fmt.Println(r.Cast1.Role+":", r.Cast1.Name)
				fmt.Println(r.Cast2.Role+":", r.Cast2.Name)
				fmt.Println()
			}
			return
		}

		path := make([]models.Relationship, len(o.Relationship.Path))
		copy(path, o.Relationship.Path)
		path = append(path, *o.Relationship)

		q.Enqueue(utilities.WorkRequest{Actor: o.Relationship.Cast2.URL, Path: path})
	}
}

func getdata(id string, res chan<- models.Response, throttle <-chan time.Time) {
	Log.Println("fetching", id)
	<-throttle
	resp, err := http.Get("http://data.moviebuff.com/" + id)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		res <- models.Response{Data: models.Data{}, Err: err}
		return
	}

	if resp.StatusCode == http.StatusForbidden {
		res <- models.Response{Data: models.Data{}, Err: models.Skippable{Reason: "forbidden"}}
		return
	}

	jsonBlob, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		res <- models.Response{Data: models.Data{}, Err: err}
		return
	}

	var data models.Data
	err = json.Unmarshal(jsonBlob, &data)
	res <- models.Response{Data: data, Err: err}
	return
}

func getrelationship(q *utilities.Queue, output chan<- models.Out) {
	visited := map[string]bool{}
	throttle := time.Tick(time.Second / 800)

	for {
		if q.Empty() {
			continue
		}

		i := q.Dequeue()
		if visited[i.Actor] {
			output <- models.Out{}
		}

		res := make(chan models.Response)
		go getdata(i.Actor, res, throttle)

		actorProfile := <-res
		if actorProfile.Err != nil {
			if _, ok := actorProfile.Err.(models.Skippable); ok {
				output <- models.Out{}
				continue
			}
			output <- models.Out{Err: actorProfile.Err}
			continue
		}

		visited[i.Actor] = true

		var count int
		urlVsMovie := map[string]models.Movie{}
		for _, Movie := range actorProfile.Data.Movies {
			if visited[Movie.URL] {
				continue
			}

			go getdata(Movie.URL, res, throttle)

			urlVsMovie[Movie.URL] = Movie
			visited[Movie.URL] = true
			count++
		}

		for k := 0; k < count; k++ {
			resp := <-res
			if resp.Err != nil {
				if _, ok := resp.Err.(models.Skippable); ok {
					continue
				}
				output <- models.Out{Err: resp.Err}
				return
			}

			for _, c := range resp.Data.Casts {
				if c.URL == i.Actor || visited[c.URL] {
					continue
				}

				visited[c.URL] = true

				output <- models.Out{
					Relationship: &models.Relationship{
						Cast1: models.Cast{Name: actorProfile.Data.Name, URL: actorProfile.Data.URL, Role: urlVsMovie[resp.Data.URL].Role},
						Cast2: c,
						Movie: urlVsMovie[resp.Data.URL].Name,
						Path:  i.Path,
					},
				}
			}

			for _, c := range resp.Data.Crew {
				if c.URL == i.Actor || visited[c.URL] {
					continue
				}

				visited[c.URL] = true
				output <- models.Out{
					Relationship: &models.Relationship{
						Cast1: models.Cast{Name: actorProfile.Data.Name, URL: actorProfile.Data.URL, Role: urlVsMovie[resp.Data.URL].Role},
						Cast2: c,
						Movie: urlVsMovie[resp.Data.URL].Name,
						Path:  i.Path,
					},
				}
			}
		}
	}
}