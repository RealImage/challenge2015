package datastructs

type General struct {
	Url  string
	Name string
}

type Entity struct {
	General
	Type string
}
type Info struct {
	General
	Role string
}

type Movie struct {
	Entity
	Cast []Info
}

type Person struct {
	Entity
	Movies []Info
}
