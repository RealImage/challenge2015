package datatype

const DataUri = "http://data.moviebuff.com/"
// moviebuff json response struct
type MoviebuffRes struct {
	Name   string  `json:"name"`
	Url    string  `json:"url"`
	Type   string  `json:"type"`
	Movies []Movie `json:"movies"`
	Cast   []CastAndCrew  `json:"cast"`
	Crew   []CastAndCrew  `json:"crew"`
}

// degree of separation structure hold to vist person & link details
type DegreesOfSeparation struct {
	Source        string
	Destination   string
	Actor1       *MoviebuffRes
	Actor2       *MoviebuffRes
	A2Movies      map[string]Movie
	VisitedPerson map[string]bool
	Visit         []string
	Visited       map[string]bool
	Link          map[string]Result
}

type Movie struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	Role string `json:"role"`
}

type CastAndCrew struct {
	Url  string `json:"url"`
	Name string `json:"name"`	
	Role string `json:"role"`
}

type Result struct {
	Movie string
	Actor1 string
	Role1 string
	Actor2 string
	Role2 string
}



