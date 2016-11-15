package dosengine

import(
	"degreesOfSeparation/httpget"
	moviebuffDatatype "degreesOfSeparation/datatype"
	"log"
	"fmt"
	"strings"
)

/**
 * [initMovieBuffData description]
 * @param  {[actor1, actor2]} ) (error [description]
 * @return {[error]}   [initialize actor1 and actor2 movie buff data]
 */
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
	return nil
}

/**
 * [findDos description]
 * @param  {[type]} ) ([]dos,       error [description]
 * @return {[type]}   [description]
 */
func DegreesOfSeparation() ([]moviebuffDatatype.Result, error) {
	var dosRes []moviebuffDatatype.Result
	for true {
		// fmt.Printf("Visited Person: %v, %d\n\n", moviebuff.VisitedPerson, len(moviebuff.VisitedPerson))
		for _, actor := range moviebuff.Visit {
			fmt.Printf("Current relation actor: %s\n\n", actor)
			if moviebuff.VisitedPerson[actor] {
				continue
			}
			moviebuff.VisitedPerson[actor] = true
			// fmt.Printf("Visited Person inner:  %v, %d\n\n", moviebuff.VisitedPerson, len(moviebuff.VisitedPerson))
			/**			
			 * @type fetch this->actor movie buff data
			 * actor1-> movie data of -> current actor in the moviebuff.Visit array
			 */
			actor1, err := httpget.FetchMoviebuffData(actor)			
			if err != nil {
				if strings.Contains(err.Error(), "looking for beginning of value") {
					continue
				}
				return nil, err
			}
			/**
			 * check the relation Among this actor2(global) and this->actor1 
			 */
			for _, a1movie := range actor1.Movies {
				// fmt.Printf("%s==%s\n\n", a1movie.Url, moviebuff.A2Movies[a1movie.Url].Url)
				if moviebuff.A2Movies[a1movie.Url].Url == a1movie.Url {
					dos := moviebuffDatatype.Result{
							a1movie.Name, 
							actor1.Name, 
							a1movie.Role, 
							moviebuff.Actor2.Name, 
							moviebuff.A2Movies[a1movie.Url].Role,
						}
					// check this actor is link bewtween them
					if _, isLinkPerson := moviebuff.Link[actor1.Url]; isLinkPerson {
						/*fmt.Printf("actor 1 name: %s\n\n",actor1.Name)
						fmt.Printf("A2MoviesName %s\n\n",moviebuff.A2Movies[a1movie.Url].Name)
						fmt.Printf("actor1Url %s\n\n",actor1.Url)
						fmt.Printf("found link person%s==%s\n\n", moviebuff.Link[actor1.Url], actor1.Url)*/
						// store result array
						dosRes = append(dosRes, moviebuff.Link[actor1.Url], dos)
					} else {
						dosRes = append(dosRes, dos)
					}
					fmt.Printf("Matches:%v\n\n",dosRes)
					return dosRes, nil
				}
			}
			/**
			 * [if one-one realtion not found]
			 * get all the crew and costing members from each of this->actor1 flim
			 * make moviebuff.Visit['actors_name']
			 */
			for _, a1movie := range actor1.Movies {
				// ignore already visted movie
				if moviebuff.Visited[a1movie.Url] {
					continue
				}
				moviebuff.Visited[a1movie.Url] = true
				/**
				 * [get this->movie.Url details]
				 * @type {[json]}
				 */
				a1moviedetail, err := httpget.FetchMoviebuffData(a1movie.Url)
				if err != nil {
					if strings.Contains(err.Error(), "looking for beginning of value") {
						continue
					}
					return nil, err
				}
				/**
				 * [get all the casting persons details]
				 * @type {[type]}
				 */
				for _, a1moviecast := range a1moviedetail.Cast {
					if moviebuff.Visited[a1moviecast.Url] {
						continue
					}
					moviebuff.Visited[a1moviecast.Url] = true
					// add this->cast actor url in Visit[array]
					moviebuff.Visit = append(moviebuff.Visit, a1moviecast.Url)
					dos := moviebuffDatatype.Result{
						a1movie.Name, 
						actor1.Name, 
						a1movie.Role, 
						a1moviecast.Name, 
						a1moviecast.Role,
					}
					moviebuff.Link[a1moviecast.Url] = dos
				}
				/**
				 * [get all the crew persons details]
				 * @type {[type]}
				 */
				for _, a1moviecrew := range a1moviedetail.Crew {
					if moviebuff.Visited[a1moviecrew.Url] {
						continue
					}
					moviebuff.Visited[a1moviecrew.Url] = true
					// add this->crew technician url in Visit[array]
					moviebuff.Visit = append(moviebuff.Visit, a1moviecrew.Url)
					dos := moviebuffDatatype.Result{
						a1movie.Name, 
						actor1.Name, 
						a1movie.Role, 
						a1moviecrew.Name, 
						a1moviecrew.Role,
					}
					moviebuff.Link[a1moviecrew.Url] = dos
				}
			}

		}
		// fmt.Printf("Link: %v, %d\n\n", moviebuff.Link, len(moviebuff.Link))
		// fmt.Printf("Visit: %v, %d\n\n", moviebuff.Visit, len(moviebuff.Visit))
		// fmt.Printf("Visited: %v, %d\n\n", moviebuff.Visited, len(moviebuff.Visited))
	}												

	return dosRes, nil
}

var checkErr = func(err error) {
	if err != nil {
		log.Fatal(err)
	}
}