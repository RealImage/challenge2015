package main

/** /
import (
	"github.com/trustmaster/goflow"
)

type DegreeGraph struct {
	flow.Graph
}

func CreateDegreeGraph(d *Degree, graphIn chan result) *DegreeGraph {
	degreeGraph := new(DegreeGraph)
	degreeGraph.InitGraphState()

	personComponent := new(PersonComponent)
	personComponent.Component.Mode = flow.ComponentModeSync
	degreeGraph.Add(personComponent, "PersonComponent")

	movieComponent := new(MovieComponent)
	movieComponent.degree = d
	personComponent.Component.Mode = flow.ComponentModeSync
	degreeGraph.Add(movieComponent, "MovieComponent")

	degreeGraph.Connect("PersonComponent", "PersonComponentOp", "MovieComponent", "MovieComponentIp")
	movieComponent.MovieComponentOp = graphIn

	degreeGraph.MapInPort("In", "PersonComponent", "PersonComponentIp")
	return degreeGraph
}
/**/
