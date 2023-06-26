package main

// Struct for parsing urls
type URL struct {
	Movies []string `json:"movies"`
	Actor  []string `json:"actors"`
}

// Structs for Data handling from the URL
type CommonData struct {
	URL  string `json:"url"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type Data struct {
	URL  string `json:"url"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type ActorData struct {
	CommonData
	Movies []Data `json:"movies"`
}

type MoiveData struct {
	CommonData
	Casts []Data `json:"cast"`
}

// Structs below are used for internal data handling and transaformation
type RoleInMovie struct {
	MovieName string `json:"name"`
	Role      string `json:"role"`
}

type CastInMovie struct {
	CastName string `json:"name"`
	Role     string `json:"role"`
}

type MovieNameCastName struct {
	Movie string `json:"movie"`
	Cast1 string `json:"cast1"`
	Cast2 string `json:"cast2"`
}
