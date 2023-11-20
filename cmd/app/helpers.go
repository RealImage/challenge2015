package main

import (
	"fmt"
	"io"
	"net/http"
)

func (app *GlobalVar) doHTTPRequest(sURL string) ([]byte, error) {
	sFunctionName := "doHTTPRequest"

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
