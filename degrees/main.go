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
	"strings"
	"time"

	"gophercon/moviebuff"
)

func main() {
	//profilling

	//defer profile.Start(profile.CPUProfile).Stop()

	//tracing

	//	f, err := os.Create(time.Now().Format("2006-01-02T150405.pprof"))
	//	if err != nil {
	//		panic(err)
	//	}
	//	defer f.Close()

	//	if err := trace.Start(f); err != nil {
	//		panic(err)
	//	}
	//	defer trace.Stop()

	//check if argument is passed correctly
	if len(os.Args) != 3 {
		log.Fatal("\nUSASE :: degrees <first-person-name><space><second-person-name>\n \tExample :: degrees amitabh-bachchan robert-de-niro")
	} else {

		//retrieve the inputs
		src := strings.ToLower(os.Args[1])
		dest := strings.ToLower(os.Args[2])

		//parse configuration file
		config, err := processConfig()
		if err != nil {
			log.Fatalln(err.Error())
		}

		//initialize the connection
		var connection moviebuff.Connection
		err = connection.Initialize(src, dest, config)
		if err != nil {
			log.Fatalln(err.Error())
		}

		t1 := time.Now()
		result, err := connection.GetConnection()
		if err != nil {
			log.Fatalf("Error in finding the degree of connection between %s and %s.\n Error :: %s", src, dest, err.Error())
		}

		//print result
		printResult(result, t1)
	}
}

//parse configuration file
func processConfig() (*moviebuff.Conf, error) {
	//read config file
	data, err := ioutil.ReadFile("conf.json")
	if err != nil {
		return nil, err
	}

	//unmarshel data
	var config moviebuff.Conf
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	//set the maximum number of process to be used
	if config.NumCPU > 0 {
		runtime.GOMAXPROCS(config.NumCPU)
	} else {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
	return &config, nil
}

//printResult prints the output in desired format
func printResult(relations []moviebuff.Relation, t1 time.Time) {
	fmt.Println("Time Taken: ", time.Since(t1))
	//display the output
	fmt.Println("\nDegree of saperation: ", len(relations))
	for i, relation := range relations {
		fmt.Printf("\n%d. Movie: %s\n%s: %s\n%s: %s\n", i+1, relation.Movie, relation.Role1, relation.Person1, relation.Role2, relation.Person2)
	}
}
