package model

type Processor interface {
	GetProcessedData(firstPerson, secondPerson string)
}

// data structure for this Problem
type PersonMovie struct {
	Name string `json:"name" bson:"name,omitempty"`
	Url  string `json:"url" bson:"url,omitempty"`
	Role string `json:"role" bson:"role,omitempty"`
}

type Person struct {
	Url  string `json:"url" bson:"url,omitempty"`
	Type string `json:"type" bson:"type,omitempty"`
	Name string `json:"name" bson:"name,omitempty"`

	Movies []PersonMovie `json:"movies" bson:"movies,omitempty"`
}

type People struct {
	Name string `json:"name" bson:"name,omitempty"`
	Url  string `json:"url" bson:"url,omitempty"`
	Role string `json:"role" bson:"role,omitempty"`
}

type Movie struct {
	Url  string `json:"url" bson:"url,omitempty"`
	Type string `json:"type" bson:"type,omitempty"`
	Name string `json:"name" bson:"name,omitempty"`

	Cast []People `json:"cast" bson:"cast,omitempty"`
	Crew []People `json:"crew" bson:"cast,omitempty"`
}

// type Crews struct {
// 	Name string `json:"name" bson:"name,omitempty"`
// 	Url  string `json:"url" bson:"url,omitempty"`
// 	Role string `json:"role" bson:"role,omitempty"`
// }

// type Result struct {
// 	MovieName string `json:"moviename" bson:"moviename,omitempty"`
// 	CastName  string `json:"castname" bson:"castname,omitempty"`
// 	Role      string `json:"role" bson:"role,omitempty"`
// }
