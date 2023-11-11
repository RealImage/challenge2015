package main

var (
	MoviesUrl = []string{"http://data.moviebuff.com/the-great-gatsby",
		"http://data.moviebuff.com/the-wolf-of-wall-street",
		"http://data.moviebuff.com/taxi-driver"}

	ActorsUrl = []string{"http://data.moviebuff.com/amitabh-bachchan",
		"http://data.moviebuff.com/leonardo-dicaprio",
		"http://data.moviebuff.com/martin-scorsese"}
	GlobalUrlActorMap = make(map[string]string)
)
