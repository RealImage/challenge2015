package dosengine
import(
	"fmt"
	"os"
	"degreesOfSeparation/httpget"
	"log"
)
func initMovieBuffData(actor1, actor2 string) error{
	act1, err := httpget.FetchMoviebuffData(actor1)
	checkErr(err)
	act2, err := httpget.FetchMoviebuffData(actor2)
	checkErr(err)
	/**
	 * [len description]
	 * @param  {[type]} act1.Movies) >             len(act2.Movies [description]
	 * @set actor 1 has less movie data
	 */
	if len(act1.Movies) > len(act2.Movies) {
		moviebuff.Source, moviebuff.Destination = actor2, actor1
		moviebuff.Actor1, moviebuff.Actor2 = act2, act1
	} else {
		moviebuff.Source, moviebuff.Destination = actor1, actor2
		moviebuff.Actor1, moviebuff.Actor2 = act1, act2
	}
	/**
	 * [movie description]
	 * Actor1 -> less movie data
	 * Actor2 -> greater movie data than actor1
	 * A2Movies -> Actor2 movie map[url] data
	 * @type {[type]}
	 */
	for _, movie := range moviebuff.Actor2.Movies {
		moviebuff.A2Movies[movie.Url] = movie
	}
	moviebuff.Visit = append(moviebuff.Visit, moviebuff.Source)
	moviebuff.Visited[moviebuff.Source] = true


	fmt.Println(httpget.TotalRequest)	
	fmt.Println("\n\n\n\n A2Movies buffer data \n\n")
	fmt.Printf("%v\n\n", moviebuff.A2Movies)
	fmt.Printf("%v\n\n", moviebuff.Actor1)
	fmt.Printf("%v\n\n", moviebuff.Source)
	os.Exit(1)
	return nil
}

var checkErr = func(err error) {
	if err != nil {
		log.Fatal(err)
	}
}