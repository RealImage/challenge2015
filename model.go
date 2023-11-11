package main

type URL struct {
	Movies []string `json:"movies"`
	Actor  []string `json:"actors"`
}

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

type ActorsData struct {
	CommonData
	Movies []Data `json:"movies"`
}

type MoivesData struct {
	CommonData
	Casts []Data `json:"cast"`
}

type RoleInMovie struct {
	MovieName string `json:"name"`
	Role      string `json:"role"`
}

type CastInMovie struct {
	CastName string `json:"name"`
	Role     string `json:"role"`
}

type MovieNameWithCastName struct {
	Movie string `json:"movie"`
	Cast1 string `json:"cast1"`
	Cast2 string `json:"cast2"`
}

type ActorMovieRelation struct {
	Actor       string
	commonMovie string
}
