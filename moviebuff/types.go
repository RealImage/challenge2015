/*
Purpose 	  : this file contains the struct used
				in parsing the jsons
File Name	  : main.go
Package		  : moviebuff
Date 		  : 01.08.2015
Author 		  : Mayank Patel
Date		Name		Modification
*/

// moviebuff project moviebuff.go
//this project get the degree of connection between
//two celebrity and tells how they are connected
package moviebuff

//conf is used to parse the configuration type
type Conf struct {
	NumCPU     int    `json:"cpu_core"`
	Address    string `json:"bucket_address"`
	Limit      int    `json:"rate-limit"`
	RetryCount int    `json:"connection-retry-count"`
}

//job store the job done by a person in a movie
type credit struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	Role string `json:"role"`
}

//url detail contain the information of one entity
//(movie or actor)
type details struct {
	Url    string   `json:"url"`
	Typ    string   `json:"type"`
	Name   string   `json:"name"`
	Movies []credit `json:"movies"`
	Cast   []credit `json:"cast"`
	Crew   []credit `json:"crew"`
}

type person struct {
	url      string
	relation []Relation
}

//Relation describe how two person are connected
//to each other
type Relation struct {
	Movie   string
	Person1 string
	Role1   string
	Person2 string
	Role2   string
}
