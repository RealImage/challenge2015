package main

type InfoNode struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Role string `json:"role"`
}

type InfoNodeForQueue struct {
	ParentNodeEntry *InfoNodeForQueue
	InfoNodeEntry   InfoNode
}

type MovieBuffResponseDataForPerson struct {
	URL    string     `json:"url"`
	Type   string     `json:"type"`
	Name   string     `json:"name"`
	Movies []InfoNode `json:"movies"`
}

type MovieBuffResponseDataForMovies struct {
	URL  string     `json:"url"`
	Type string     `json:"type"`
	Name string     `json:"name"`
	Cast []InfoNode `json:"cast"`
	Crew []InfoNode `json:"crew"`
}
