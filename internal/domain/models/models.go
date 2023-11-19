package moviebuff

type Data struct {
	Cast   []Data `json:"cast"`
	Crew   []Data `json:"crew"`
	Movies []Data `json:"movies"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	Type   string `json:"type"`
	Url    string `json:"url"`
}

type ActorData struct {
	Url    string `json:"url"`
	Type   string `json:"type"`
	Name   string `json:"name"`
	Movies []Data `json:"movies"`
}

type MovieData struct {
	Url  string `json:"url"`
	Type string `json:"type"`
	Cast []Data `json:"cast"`
	Crew []Data `json:"crew"`
}
