package dtos

type ActorDetails struct {
	URL            string         `json:"url"`
	Type           string         `json:"type"`
	Name           string         `json:"name"`
	MoviesAndRoles []MovieAndRole `json:"movies"`
}

type MovieAndRole struct {
	URL  string `json:"url"` //movie's url
	Name string `json:"name"`
	Role string `json:"role"`
}

type MovieDetail struct {
	URL  string        `json:"url"`
	Type string        `json:"type"`
	Name string        `json:"name"`
	Cast []CastAndCrew `json:"cast"`
	Crew []CastAndCrew `json:"crew"`
}

type CastAndCrew struct {
	URL  string `json:"url"` //person's URL
	Name string `json:"name"`
	Role string `json:"role"`
}
