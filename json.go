package main

import (
	"encoding/json"	
	"io/ioutil"
	"net/http"
)

type movie struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	Role string `json:"role"`
}

type cast struct {
	Url  string `json:"url"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type crew struct{
	Name string `json:"name"`
	Url  string `json:"url"`
	Role string `json:"role"`
}

type cons struct{
	Name   string  `json:"name"`
	Url    string  `json:"url"`
	Type   string  `json:"type"`
	Movies []movie `json:"movies"`
	Cast   []cast  `json:"cast"`
	Crew   []crew  `json:"crew"`
}

func getData(url string) (*cons, error) {
	resp, err := http.Get(url)
	defer ErrHandle(err)
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	var data cons
	err = json.Unmarshal(body, &data)
	return &data,  nil
}

