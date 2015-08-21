//Degrees project main.go
//this project get the degree of connection between
//two celebrity and tells how they are connected
package main

import (
	"gophercon/moviebuff"
	"log"
	"testing"
	"time"
)

func Test(t *testing.T) {
	//parse configuration file
	config, err := processConfig()
	if err != nil {
		log.Fatalln(err.Error())
	}

	//initialize the connection
	var connection moviebuff.Connection
	err = connection.Initialize("winfield-scott-mattraw", "welker-white", config)
	if err != nil {
		log.Fatalln(err.Error())
	}

	t1 := time.Now()
	result, err := connection.GetConnection()
	if err != nil {
		log.Fatalf("Error in finding the degree of connection.\n Error :: %s", err.Error())
	}

	//print result
	printResult(result, t1)
}
