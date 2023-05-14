package utility

import (
	"encoding/json"
	"fmt"

	httpClient "qubecinemas.com/clients"
	model "qubecinemas.com/internal/models"
)

//repeated const values
const (
	baseURL = "http://data.moviebuff.com"
)

var client = httpClient.NewHTTPClient(baseURL)

func PersonUtility(personName string) (error, *model.Person) {

	response, _, err := client.MakeRequest("GET", fmt.Sprintf("%s/%s", client.BaseURL, personName), nil)
	if err != nil {
		return err, nil
	}

	var person model.Person
	json.Unmarshal(response, &person)

	return nil, &person
}

func MovieUtility(movieName string) (error, *model.Movie) {
	response, _, err := client.MakeRequest("GET", fmt.Sprintf("%s/%s", client.BaseURL, movieName), nil)
	if err != nil {
		return err, nil
	}
	var movieSub model.Movie
	json.Unmarshal(response, &movieSub)

	return nil, &movieSub
}
