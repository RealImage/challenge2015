package main

type PersonMovies struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	Role string `json:"role"`
}

type Person struct {
	Url    string         `json:"url"`
	Type   string         `json:"type"`
	Name   string         `json:"name"`
	Movies []PersonMovies `json:"movies"`
}

type MovieCast struct {
	Url  string `json:"url"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type MovieCrew struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	Role string `json:"role"`
}

type Movie struct {
	Url  string      `json:"url"`
	Type string      `json:"type"`
	Name string      `json:"name"`
	Cast []MovieCast `json:"cast"`
	Crew []MovieCrew `json:"crew"`
}

type match struct {
	key       string
	srcValue  string
	destValue string
}
