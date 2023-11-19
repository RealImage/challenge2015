package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func main() {
	args := os.Args
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("\nPanic occurred: %+v", err)
			return
		}
	}()
	if len(args) < 3 || len(args) > 3 {
		fmt.Println("invalid command line")
	}
	cmd := args[1]
	var degree int

	switch cmd {
	default:
		fmt.Println("\nInvalid command command line")
		return
	case "degrees":
		fmt.Scanf("Degrees of Separation: %d", &degree)
	}

	movie := args[2]
	person := args[3]

	fetchInfo(movie, person)
}

func fetchInfo(movie, person string) {

}

type Movie struct {
	URL  string `json:"url"`
	Name string `json:"name"`
	Type string `json:"type"`
	Cast []cast `json:"cast"`
}

// type cast struct {
// 	Name string `json:"name"`
// 	URL  string `json:"url"`
// 	Role string `json:"role"`
// }

type Person struct {
	URL    string `json:"url"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Movies []cast `json:"movies"`
}

func fetchMovie(movie, person string) {

	baseURL := "http://data.moviebuff.com/"
	movieURL := fmt.Sprint("%s%s", baseURL, movie)
	req, err := http.NewRequest(http.MethodGet, movieURL, nil)
	if err != nil {
		//handle error
	}

	tlsConfig := tls.Config{
		MinVersion:         tls.VersionTLS13,
		InsecureSkipVerify: true,
	}
	// Skipping verification as it is done manually using digest of the TLS certificate as this is step of setting up service
	transport := http.Transport{
		TLSClientConfig: &tlsConfig,
	}
	client := &http.Client{Transport: &transport}
	resp, err := client.Do(req)
	if err != nil {

	}

	err = json.NewDecoder(resp.Body).Decode(&userCreateResponse)
	if err != nil {
		return nil, err
	}

}
