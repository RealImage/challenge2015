package main

import "encoding/json"

type header struct {
	Url  string `json:"url"`
	Typ  string `json:"Type"`
	Name string `json:"name"`
}

type movielist struct {
	header
	Movies []movie `json: "movies"`
}

type movie struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	Role string `json:"role"`
}

type castList struct {
	header
	CastL []cast `json:"cast"`
}

type cast struct {
	Url  string `json:"url"`
	Name string `json:"name"`
	Role string `json:"role"`
}

func getMovies(body []byte) ([]movie, error) {
	var movies movielist
	err := json.Unmarshal(body, &movies)
	ErrHandle(err)
	return movies.Movies, nil
}

func getCast(body []byte) ([]cast, error) {
	var persons castList
	err := json.Unmarshal(body, &persons)
	ErrHandle(err)
	return persons.CastL, nil
}
