package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (app *GlobalVar) doHTTPRequest(sURL string) ([]byte, error) {
	sFunctionName := "doHTTPRequest"

	if !strings.Contains(sURL, "amitabh-bachchan") &&
		!strings.Contains(sURL, "the-great-gatsby") &&
		!strings.Contains(sURL, "leonardo-dicaprio") &&
		!strings.Contains(sURL, "the-wolf-of-wall-street") &&
		!strings.Contains(sURL, "martin-scorsese") &&
		!strings.Contains(sURL, "taxi-driver") {
		return nil, errors.New("Invalid value")
	}

	// Make GET request
	response, err := http.Get(sURL)
	if err != nil {
		fmt.Println(sFunctionName, "Error making the request:", err)
		return nil, err
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading the response body:", err)
		return nil, err
	}

	return body, nil
}
