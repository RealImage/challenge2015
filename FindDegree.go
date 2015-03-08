package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"sync"
	"time"
)

const url = "http://data.moviebuff.com/"

type jsonMovie struct {
	Url  string  `json:"url"`
	Type string  `json:"type"`
	Name string  `json:"name"`
	Cast []Movie `json:"cast"`
	Crew []Movie `json:"crew"`
}

type Movie struct {
	Name string
	Url  string
	Role string
}

type Person struct {
	Name      string  `jason:"name"`
	Url       string  `jason:"url"`
	Type      string  `jason:"type"`
	MovieList []Movie `json:"movies"` //using map very expensive
}

type VMovie struct {
	movieMap map[string]bool // we need to sync this!
	sync.RWMutex
}

type VPerson struct {
	personMap map[string]bool // we need to sync this!
	sync.RWMutex
}

var VisitedMovie *VMovie
var VisitedPerson *VPerson

func (v *VPerson) checkVisitedPerson(url string) (val, ok bool) {

	v.RLock()
	defer v.RUnlock()
	val, ok = v.personMap[url]
	return

}

func (v *VMovie) checkVisitedMovie(url string) (val, ok bool) {

	v.RLock()
	defer v.RUnlock()
	val, ok = v.movieMap[url]
	return

}

func (v *VMovie) MarkVisitedMovie(url string) {
	v.Lock()
	defer v.Unlock()
	v.movieMap[url] = true
}

func (v *VPerson) MarkVisitedPerson(url string) {
	v.Lock()
	defer v.Unlock()
	v.personMap[url] = true

}

// ******************************************************

/*
This will hold the relation graph between persons
*/

var gPersonStackMutex *sync.Mutex = &sync.Mutex{}

type PersonStack struct {
	Person    string // from the current person
	PersonUrl string //from the current person
	MovieName string //from cast n crew
	MovieUrl  string //from cast n crew
	Role      string // from cast n crew

}

type gPersonStack struct {
	p []PersonStack
}

var globalPersonStack []gPersonStack //array of PersonStack path in a stack

func printResult(p []PersonStack) {
	j := 1
	p = p[1:]

	fmt.Print("\nDegrees of Separation: ", (len(p))/2, "\n")
	for i := 0; i < len(p); i = i + 2 {
		fmt.Print(j, ". Movie : ", p[i].MovieName, "\n")
		fmt.Print(p[i].Role, " : ", p[i].Person, "\n")
		fmt.Print(p[i+1].Role, " : ", p[i+1].Person, "\n\n")
		j++
	}
}

/*
* this function builds the PersonStack
* we need error control in this function say for some reason the passed object are not proper
* TODO
 */

func getBuildtheStackObject(perObj *Person, movie *Movie) PersonStack {

	obj := PersonStack{Person: perObj.Name, PersonUrl: perObj.Url, MovieName: movie.Name, MovieUrl: movie.Url, Role: movie.Role}
	return obj
}

func initglobals() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	globalPersonStack = make([]gPersonStack, 0)
	VisitedMovie = &VMovie{movieMap: make(map[string]bool)}
	VisitedPerson = &VPerson{personMap: make(map[string]bool)}
}

func putgPersonStack(p *gPersonStack) {

	gPersonStackMutex.Lock()

	globalPersonStack = append(globalPersonStack, *p)

	gPersonStackMutex.Unlock()

}

func getgPersonStack() (p gPersonStack, err bool) {

	gPersonStackMutex.Lock()
	lenn := len(globalPersonStack)
	gPersonStackMutex.Unlock()
	if lenn <= 0 {
		err = true
		p = gPersonStack{}
	} else {

		gPersonStackMutex.Lock()
		res := globalPersonStack[0]
		globalPersonStack = globalPersonStack[1:lenn]
		gPersonStackMutex.Unlock()
		err = false
		p = res

	}

	return

}

/*
this shld be the last person in the build tree
*/
func getNextPersonFromPersonStack(pp []PersonStack) *PersonStack {
	if len(pp) > 0 {
		return &pp[len(pp)-1]
	}

	return nil
}

/*
This function will have to fetch the url passed and also maintain the
rate limit...
*/

func getPersonDetail(urlname string, sendPerson chan Person, errorchan chan bool) {

	var actor Person
	var flagAtemps int = 0
	for {

		if urlname == "" {
			fmt.Println("Url is null")
			errorchan <- true
			return
		}
		resp, err := http.Get((url + urlname))
		if err != nil || (resp != nil && resp.StatusCode != 200) {

			flagAtemps++
			if flagAtemps <= 500 {
				time.Sleep(300 * time.Millisecond)
				continue
			} else {
				fmt.Println("\n Failed to fetch person", urlname)
				errorchan <- true
				return
			}

		}

		defer resp.Body.Close()
		jasontxt, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic("io error")
			errorchan <- true
			return
		}

		err = json.Unmarshal(jasontxt, &actor)

		if err != nil {
			errorchan <- true
			return
		}

		if len(actor.MovieList) == 0 {
			errorchan <- true
			return
		}

		sendPerson <- actor
		break

	}

	return

}

/*
* This function tries to get the Movie data.
* and calls the getUniqueCastnCrew
 */

func getCastCrewnMovie(prefixStack []PersonStack, urlmovie Movie, person Person,
	donechan chan bool, stopchan chan bool, stopchanLen int, entityTwo string) {

	var actor jsonMovie
	var ok bool
	var flagAtemps int = 0

	defer VisitedMovie.MarkVisitedMovie(urlmovie.Url) //even if we fail to fetch we mark it as done

	for {
		resp, err := http.Get((url + urlmovie.Url))

		if err != nil || (resp != nil && resp.StatusCode != 200) {

			flagAtemps++

			if flagAtemps <= 500 {

				time.Sleep(300 * time.Millisecond)

				continue
			} else {

				fmt.Println("\nFailed to fetch the Movie details", urlmovie.Url)

				donechan <- false
				return
			}

		}

		defer resp.Body.Close()

		jasontxt, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic("io error ")
			donechan <- false
			return
		}
		err = json.Unmarshal(jasontxt, &actor)

		if err != nil {
			donechan <- false
			return
		}

		ok = getUniqueCastnCrew(prefixStack, person, &urlmovie, &actor, stopchan, stopchanLen, entityTwo)
		break

	}

	donechan <- ok

	return

}

/*
* Spawns a number of goroutines to get the movies list
*
 */
func getRateLimitedCastnCrew(prefixStack []PersonStack, person Person, done chan bool, entityTwo string) {

	var stopSignal bool = false

	movielen := len(person.MovieList)

	/*
	* what is this case?? Do we need this ??
	 */

	if movielen == 0 {

		done <- true
		return
	}

	//This is a simple rate limiter need to implement a busty one later TODO

	rateLimit := time.Tick(10 * time.Millisecond)

	doneChan := make(chan bool, (movielen - 1))
	stopChan := make(chan bool, (movielen - 1)) //this is for cancellation of all goroutine

	var visitedCnt int = 0

	for _, movie := range person.MovieList {

		if _, ok := VisitedMovie.checkVisitedMovie(movie.Url); ok {
			visitedCnt++
			continue
		}

		<-rateLimit
		go getCastCrewnMovie(prefixStack, movie, person, doneChan, stopChan, movielen, entityTwo)

	}

	for i := 0; i < ((movielen) - visitedCnt); i++ {
		select {

		case <-doneChan:
		case <-stopChan:
			{

				stopSignal = true

			}

		}
	}

	done <- !stopSignal

	return

}

/*
* finds unique person and populated the next possible paths in the tree
 */

func getUniqueCastnCrew(prefixStack []PersonStack, person Person, movie *Movie, result *jsonMovie, stopchan chan bool, stopchanLen int, entityTwo string) bool {

	localPersonMap := make(map[string]bool)

	res := append(result.Cast, result.Crew...)

	newPersonStackSlice := make([]PersonStack, 0)
	prefix := getBuildtheStackObject(&person, movie)
	newPersonStackSlice = append(newPersonStackSlice, prefixStack...)
	newPersonStackSlice = append(newPersonStackSlice, prefix)

	for _, val := range res {
		_, okay := VisitedPerson.checkVisitedPerson(val.Url)
		_, ok := localPersonMap[val.Url]

		if !ok && !okay {

			localPersonMap[val.Url] = true

			if val.Role == "Production Company / Production" || val.Role == "Distributor / Distribution" || val.Role == "Music Label / Music" || val.Role == "Associate Production Company / Production" || val.Role == "Visual Effects Studio / Visual Effects" || val.Role == "Special Effects Studio / Special Effects" {
				continue
			}

			tempPersonStack := append(newPersonStackSlice, PersonStack{Person: val.Name, PersonUrl: val.Url, MovieName: result.Name, MovieUrl: result.Url, Role: val.Role})

			if val.Url == entityTwo {
				printResult(tempPersonStack)
				stopchan <- true
				return true
			}

			putgPersonStack(&gPersonStack{p: tempPersonStack})

		}

	}

	return true

}

func populateInitialStack(entityOne ...string) (PersonTwo string, stop bool) {

	stop = false
	xx := make([]Person, 2)
	var result Person

	tempMap0 := make(map[string]Movie)
	tempMap1 := make(map[string]Movie)

	perosnChan := make(chan Person, 2)
	errChan := make(chan bool, 2)

	for _, val := range entityOne {
		go getPersonDetail(val, perosnChan, errChan)

	}

	for i := 0; i < 2; i++ {
		select {
		case yy := <-perosnChan:
			{
				xx[i] = yy

				for _, val := range xx[i].MovieList {

					if i == 0 {
						tempMap0[val.Url] = val
					} else {
						tempMap1[val.Url] = val
					}
				}

			}
		case _ = <-errChan:
			{
				//fmt.Println("err received ")
			}
		}

	}

	for key1, val1 := range tempMap0 {

		if val2, ok := tempMap1[key1]; ok {
			fmt.Print("\nDegrees of separation : 1\n", "1. Movie: ", val1.Name, "\n")
			fmt.Print(val1.Role, " : ", xx[0].Name, "\n", val2.Role, " : ", xx[1].Name)
			stop = true
			return
		}

	}

	len1 := len(xx[0].MovieList)
	len2 := len(xx[1].MovieList)

	if len1 == 0 {
		fmt.Println("This Url ", xx[0].Url, " Doesn't exist so aborting further search")
		stop = true
		return

	} else if len1 == 0 {
		fmt.Println("This Url ", xx[1].Url, " Doesn't exist so aborting further search")
		stop = true
		return
	}

	if (len1 > len2) && (len2 > 0) {
		result = xx[1]
		PersonTwo = xx[0].Url

	} else {
		result = xx[0]
		PersonTwo = xx[1].Url

	}

	tempPer := make([]PersonStack, 0)
	tempPer = append(tempPer, PersonStack{Person: result.Name, PersonUrl: result.Url})

	putgPersonStack(&gPersonStack{p: tempPer})

	return

}

func getRelationDegree(entityOne string, PersonTwo string) {

	localPersonChan := make(chan Person) //unbuffered chan why? cos we pick person data sequntially
	localErroChan := make(chan bool)
	done := make(chan bool, 1)
	var err, stop bool
	var currPerson gPersonStack
	var entityTwo string

	if entityTwo, stop = populateInitialStack(entityOne, PersonTwo); stop {
		return
	}

	for {
		if currPerson, err = getgPersonStack(); err {
			//no match
			fmt.Println("All done no match found")
			return
		}

		nextPerson := getNextPersonFromPersonStack(currPerson.p)
		if _, ok := VisitedPerson.checkVisitedPerson(nextPerson.PersonUrl); ok {
			continue
		}

		go getPersonDetail(nextPerson.PersonUrl, localPersonChan, localErroChan)

		select {
		case p := <-localPersonChan:
			{

				go getRateLimitedCastnCrew(currPerson.p, p, done, entityTwo)
			}
		case <-localErroChan:
			{

				continue
			}
		}

		okay := <-done
		if okay == false {

			break
		}

		VisitedPerson.MarkVisitedPerson(nextPerson.PersonUrl)

	}

	return

}

func main() {

	flag.Parse()
	args := flag.Args()
	personOne, personTwo := args[0], args[1]

	initglobals()
	getRelationDegree(personOne, personTwo)

}
