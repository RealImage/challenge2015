/*
Purpose 	  : this file contain the maon function
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

//conf is used to parse the configuration type
type conf struct {
	NumCPU  int    `json:"cpu_core"`
	Address string `json:"bucket_address"`
}

//job store the job done by a person in a movie
type credit struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	Role string `json:"role"`
}

//url detail contain the information of one entity
//(movie or actor)
type Details struct {
	Url    string   `json:"url"`
	Typ    string   `json:"type"`
	Name   string   `json:"name"`
	Movies []credit `json:"movies"`
	Cast   []credit `json:"cast"`
	Crew   []credit `json:"crew"`
}

type url struct {
	url      string
	relation []relation
}

//relation describe how two person are connected
//to each other
type relation struct {
	movie   string
	person1 string
	role1   string
	person2 string
	role2   string
}
