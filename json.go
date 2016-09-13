package main

import "encoding/json"

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

type cons struct{
	Name   string  `json:"name"`
	Url    string  `json:"url"`
	Type   string  `json:"type"`
	Movies []movie `json:"movies"`
	Cast   []cast  `json:"cast"`
}

func getData(body []byte) (*cons, error) {
	var data cons
	err := json.Unmarshal(body, &data)
	ErrHandle(err)
	return &data,  nil
}

