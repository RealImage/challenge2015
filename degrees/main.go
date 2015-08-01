/*
Purpose 	  : this file contain the main function
				to get degree of connection
File Name	  : main.go
Package		  : main
Date 		  : 01.08.2015
Author 		  : Mayank Patel
Date		Name		Modification
*/

//Degrees project main.go
//this project get the degree of connection between
//two celebrity and tells how they are connected
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"time"
)

func main() {
	profile.

	//check if argument is passed correctly
	if len(os.Args) != 3 {
		log.Fatal("\nUSASE :: degrees <first-person-name><space><second-person-name>\n \tExample :: degrees amitabh-bachchan robert-de-niro")
	} else {

		//retrieve the inputs
		src := os.Args[1]
		dest := os.Args[2]

		//parse configuration file
		config, err := processConfig()
		if err != nil {
			log.Fatalln(err.Error())
		}

		//initialize the connection
		var connection Connection
		err = connection.Initialize(src, dest, config.Address)
		if err != nil {
			log.Fatalln(err.Error())
		}

		t1 := time.Now()
		//Get relationship
		relations, err := connection.GetRelationship()
		if err != nil {
			log.Fatalf("Error in finding the degree of connection between %s and %s.\n Error :: %s", src, dest, err.Error())
		}
		fmt.Println("Time Taken: ", time.Since(t1))

		printResult(relations)
	}
}

//parse configuration file
func processConfig() (*conf, error) {
	//read config file
	data, err := ioutil.ReadFile("conf.json")
	if err != nil {
		return nil, err
	}

	//unmarshel data
	var config conf
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	if config.NumCPU != 0 {
		runtime.GOMAXPROCS(config.NumCPU)
	} else {
		runtime.GOMAXPROCS(runtime.NumCPU() / 2)
	}
	return &config, nil
}

//printResult prints the output in desired format
func printResult(relations []relation) {
	//display the output
	fmt.Println("\nDegree of saperation: ", len(relations))
	for i, relation := range relations {
		fmt.Printf("\n%d. Movie: %s\n%s: %s\n%s: %s\n", i+1, relation.movie, relation.role1, relation.person1, relation.role2, relation.person2)
	}
}
