package main

/** /
import (
	"fmt"
	"github.com/trustmaster/goflow"
)

// Component which deals with person
type PersonComponent struct {
	flow.Component
	PersonComponentIp chan result
	PersonComponentOp chan<- result
}

func (p *PersonComponent) OnPersonComponentIp(r result) {
	if r.currentUrl == "leonardo-dicaprio" {
		fmt.Println(r.currentUrl)
	}
	var err error
	var per person
	fmt.Println(r.currentUrl)
	for i := 0; i < 20; i++ {
		per, err = getPersonData(r.currentUrl)
		if err == nil {
			break
		}
	}
	if err != nil {
		close(p.PersonComponentIp)
	} else {
		for i := range per.Movies {
			//connections := r.connections
			//currentConn := connection{
			//	person.Movies[i],
			//	associate{person.Url, person.Name, person.Movies[i].Role},
			//	associate{}}
			//res := result{person.Movies[i].Url, r.currentDegree, currentConn, connections}
			if per.Movies[i].Url != r.prevUrl {
				res := result{per.Movies[i].Url, r.currentDegree, r.currentUrl}
				p.PersonComponentOp <- res
			}
		}
	}
}

// Component which deals with movies
type MovieComponent struct {
	flow.Component
	degree           *Degree
	MovieComponentIp <-chan result
	MovieComponentOp chan result
}

func (p *MovieComponent) OnMovieComponentIp(r result) {
	if r.currentUrl == "the-great-gatsby" || r.currentUrl == "the-wolf-of-wall-street" {
		fmt.Println(r.currentUrl)
	}
	var err error
	var m movie
	fmt.Println(r.currentUrl)
	for i := 0; i < 20; i++ {
		m, err = getMovieData(r.currentUrl)
		if err == nil {
			break
		}
	}
	if err != nil {
		close(p.MovieComponentOp)
	} else {
		r.currentDegree++
		for i := range m.Cast {
			if r.currentDegree < p.degree.degree || p.degree.degree == 0 {
				if m.Cast[i].Url != r.prevUrl {
					//connections := r.connections
					//r.currentConnection.second = movie.Cast[i]
					//connections = append(connections, r.currentConnection)
					//res := result{movie.Cast[i].Url, r.currentDegree, r.currentConnection, connections}
					res := result{m.Cast[i].Url, r.currentDegree, r.currentUrl}
					if m.Cast[i].Name == p.degree.target {
						p.degree.degree = r.currentDegree
						p.degree.res = res
						close(p.MovieComponentOp)
						break
					} else {
						//fmt.Println(res)
						p.MovieComponentOp <- res
					}
				}
			} else {
				break
			}
		}
		for i := range m.Crew {
			if r.currentDegree < p.degree.degree || p.degree.degree == 0 {
				if m.Crew[i].Url != r.prevUrl {
					//connections := r.connections
					//r.currentConnection.second = movie.Crew[i]
					//connections = append(connections, r.currentConnection)
					//res := result{movie.Crew[i].Url, r.currentDegree, r.currentConnection, connections}
					res := result{m.Crew[i].Url, r.currentDegree, r.currentUrl}
					if m.Crew[i].Name == p.degree.target {
						p.degree.degree = r.currentDegree
						p.degree.res = res
						close(p.MovieComponentOp)
						break
					} else {
						//fmt.Println(res)
						p.MovieComponentOp <- res
					}
				}
			} else {
				break
			}
		}
	}
}
/**/
