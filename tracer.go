package main

import "fmt"

//uses the global trace map

type traceData struct {
	movie     string
	parent    string
	actorName string
	role      string
}

func (t *traceData) addTrace(movie, parent string) {
	t.movie = movie
	t.parent = parent
}

func tracer(child string) {
	//trace is a global map
	if trace[child].parent != ""{
		defer fmt.Println(child, trace[child])
		tracer(trace[child].parent)
	}
}
