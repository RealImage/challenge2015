package dosengine
import (
	"fmt"
	// "os"
	"time"
	moviebuffDatatype "degreesOfSeparation/datatype"
	request "degreesOfSeparation/httpget"
)
var moviebuff moviebuffDatatype.DegreesOfSeparation

func DoS_Result(actor1, actor2 string) string{	
	// Initialize Empty Map Variable
	moviebuff.A2Movies = make(map[string]moviebuffDatatype.Movie)		
	moviebuff.Visited  = make(map[string]bool)
	moviebuff.Link = make(map[string]moviebuffDatatype.Result)
	moviebuff.VisitedPerson = make(map[string]bool)	
	err := initMovieBuffData(actor1, actor2)
	checkErr(err)
	t1 := time.Now()
	// find relation among the two actors
	result, err := DegreesOfSeparation()
	checkErr(err)
	t2 := time.Now()

	// format result data
	fmt.Printf("\nDegree of separation: %d\n\n", len(result))
	for i, d := range result {
		fmt.Printf("%d. Movie: %s\n   %s: %s\n   %s: %s\n\n", i+1, d.Movie, d.Role1, d.Actor1, d.Role2, d.Actor2)
	}

	// Optional stats
	fmt.Println("Total HTTP request sent: ", request.TotalRequest)
	fmt.Println("Time taken: ", t2.Sub(t1))

	
	return "DOS"
}