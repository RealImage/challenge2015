package dosengine

import (
	"bytes"
	"fmt"
	"html/template"
	"time"
	moviebuffDatatype "degreesOfSeparation/datatype"
	request "degreesOfSeparation/httpget"
)
var moviebuff moviebuffDatatype.DegreesOfSeparation
const ResultTemplate = `
Degree of separation: {{.DoS}}
<br />Relation Among Two Actors :<br />
{{ range $key, $value := .Relation }}
   <strong>{{ inc $key }} Movie</strong>: {{ $value.Movie }}<br />
   <strong>{{ $value.Role1 }} </strong>: {{ $value.Actor1 }}<br />
   <strong>{{ $value.Role2 }} </strong>: {{ $value.Actor2 }}<br /><br />
{{ end }}
Total HTTP request sent http://data.moviebuff.com/ : {{.HttpReq}}
<br />Time taken to find Degrees of separation: {{.Time}}`				

func DoS_Result(actor1, actor2 string) string{	
	// Initialize Empty Map Variable, clear the previous request slice data's
	moviebuff.A2Movies = make(map[string]moviebuffDatatype.Movie)		
	moviebuff.Visited  = make(map[string]bool)
	moviebuff.Link = make(map[string]moviebuffDatatype.Result)
	moviebuff.VisitedPerson = make(map[string]bool)	
	moviebuff.Actor1 = &moviebuffDatatype.MoviebuffRes{}
	moviebuff.Actor2 = &moviebuffDatatype.MoviebuffRes{}
	moviebuff.Visit = []string{}
	/*fmt.Printf("moviebuffA2Movies%v\n", moviebuff.A2Movies)
	fmt.Printf("moviebuffVisited%v\n",moviebuff.Visited)
	fmt.Printf("moviebuffLink%v\n",moviebuff.Link)
	fmt.Printf("moviebuffVisit%v\n",moviebuff.Visit)	
	fmt.Printf("moviebuffVisitedPerson%v\n",moviebuff.VisitedPerson)
	fmt.Printf("moviebuffActor1%v\n",moviebuff.Actor1)
	fmt.Printf("moviebuffActor2%v\n\n",moviebuff.Actor2)
	fmt.Println(actor1)
	fmt.Println(actor2)
	fmt.Println()/*
	err := initMovieBuffData(actor1, actor2)
	checkErr(err)
	t1 := time.Now()
	// find relation among the two actors
	result, err := DegreesOfSeparation()
	checkErr(err)
	t2 := time.Now()
	// function map to increment $key value
	funcMap := template.FuncMap{
		// The name "inc" is what the function will be called in the template text.
		"inc": func(i int) int {
			return i + 1
		},
	}
	data := map[string]interface{}{
		"DoS": len(result),
		"HttpReq": request.TotalRequest,
		"Time": t2.Sub(t1),
		"Relation": result,
	}
	t := template.Must(template.New("result").Funcs(funcMap).Parse(ResultTemplate))
	buf := &bytes.Buffer{}
	if err := t.Execute(buf, data); err != nil {
		panic(err)
	}
	s := buf.String()
	// reset request.TotalRequest to 0
	request.TotalRequest = 0	
	return s
}