package main

import "fmt"

//uses the global trace map

type traceData struct {
	movie     string
	parent    string
	actorName string
	role      string
}

func (t *traceData) addTrace(movie, parent, actorName, role string) {
	t.movie = movie
	t.parent = parent
	t.actorName = actorName
	t.role = role
}

func tracer(child, parent string) {
	//trace is a global map
	if trace[child].parent != "" && trace[child].parent != child {
		defer prettyPrint(trace[child].movie, 
		trace[parent].role,
		trace[parent].actorName,
		trace[child].role,
		trace[child].actorName) 
		tracer(trace[child].parent, trace[parent].parent)
	}
}

func prettyPrint(movie, role1, name1, role2, name2 string){
	fmt.Println("")
	fmt.Println("Movie: ", movie)
	fmt.Println(role1, ": ", name1)
	fmt.Println(role2, ": ", name2)
}
