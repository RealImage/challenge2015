package dosengine
import (
	// "fmt"
	// "os"
	moviebuffDatatype "degreesOfSeparation/datatype"
)
var moviebuff moviebuffDatatype.DegreesOfSeparation

func DoS_Result(actor1, actor2 string) string{
	_ = initMovieBuffData(actor1, actor2)
	// Initialize Empty Map Variable
	moviebuff.A2Movies = make(map[string]moviebuffDatatype.Movie)		
	moviebuff.Visited  = make(map[string]bool)
	moviebuff.Link = make(map[string]moviebuffDatatype.Result)
	moviebuff.VisitedPerson = make(map[string]bool)	
	
	return "DOS"
}