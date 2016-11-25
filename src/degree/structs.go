package main

type result struct {
	currentUrl    []string
	currentDegree int
	//currentConnection connection
	//connections       []connection
}

type urls struct {
	url         string
	connections []connection
}

type connection struct {
	movie  associate
	first  associate
	second associate
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
