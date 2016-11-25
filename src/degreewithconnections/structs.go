package main

type result struct {
	currentUrls   []urls
	currentDegree int
}

type urls struct {
	url         string
	connections []connection
}

type connection struct {
	movie      string
	first      string
	firstRole  string
	second     string
	secondRole string
}

type associate struct {
	Url  string `json:"url"`
	Name string `json:"name"`
	Role string `json:"role"`
}
type person struct {
	Url    string      `json:"url"`
	Type   string      `json:"type"`
	Name   string      `json:"name"`
	Movies []associate `json:"movies"`
}

type movie struct {
	Url  string      `json:"url"`
	Type string      `json:"type"`
	Name string      `json:"name"`
	Cast []associate `json:"cast"`
	Crew []associate `json:"crew"`
}
