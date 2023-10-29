package models
type Cast struct {
	Name string `json:"name"`
	Role string `json:"role"`
	URL  string `json:"url"`
}

type Movie struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Role string `json:"role"`
}

type Data struct {
	URL    string  `json:"url"`
	Type   string  `json:"type"`
	Name   string  `json:"name"`
	Movies []Movie `json:"movies"`
	Casts  []Cast  `json:"cast"`
	Crew   []Cast  `json:"crew"`
}

type Out struct {
	Relationship   *Relationship
	Err error
}

type Skippable struct {
	Reason string
}

func (s Skippable) Error() string {
	return s.Reason
}

type Relationship struct {
	Cast1 Cast
	Cast2 Cast
	Movie string
	Path  []Relationship
}

type Response struct {
	Data   Data
	Err error
}