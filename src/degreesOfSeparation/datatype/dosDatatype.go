package datatype

import (
	"encoding/json"
)

const DataUri = "http://data.moviebuff.com/"
// moviebuff json response struct
type MoviebuffRes struct {
	Name   string  `json:"name"`
	Url    string  `json:"url"`
	Type   string  `json:"type"`
	Movies []Movie `json:"movies"`
	Cast   []CastAndCrew  `json:"cast"`
	Crew   []CastAndCrew  `json:"crew"`
}

// Primary struct to hold the execution data
type DegreesOfSeparation struct {
	source        string
	destination   string
	person1       *MoviebuffRes
	person2       *MoviebuffRes
	p2Movies      map[string]Movie
	visitedPerson map[string]bool
	visit         []string
	visited       map[string]bool
	link          map[string]dos
}

type Movie struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	Role string `json:"role"`
}

type CastAndCrew struct {
	Url  string `json:"url"`
	Name string `json:"name"`	
	Role string `json:"role"`
}

type Result struct {
	movie          string
	person1, role1 string
	person2, role2 string
}



