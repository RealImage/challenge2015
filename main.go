package main

import (
    "fmt"
	"os"
	"github.com/challenge2015/util"
	"github.com/challenge2015/service"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Error: Invalid number of arguments. Please provide two actor names.")
		fmt.Println("Usage: degrees <actor-name> <actor-name>")
		return
	}

	firstActorName := os.Args[1]
	secondActorName := os.Args[2]

	if !util.IsValidURL(firstActorName) || !util.IsValidURL(secondActorName) {
		fmt.Println("Error: Invalid URL format. URLs should be in the format 'actor-name' or 'director-name' separated by a hyphen.")
		return
	}

	fmt.Println("First actor name: ", firstActorName)
	fmt.Println("Second actor name: ", secondActorName)

	service.DisplayResult(firstActorName, secondActorName)
}