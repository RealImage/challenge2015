package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"sync"
	"time"
)

type Degree struct {
	source       string
	target       string
	degree       int
	res          result
	graphIn      chan result
	out          chan bool
	nextUrlsChan chan []urls
	urlsChan     chan []urls
	personUrls   map[string]bool
	movieUrls    map[string]bool
	wg           sync.WaitGroup
	nextUrls     []urls
	connections  []connection
	sync.Mutex
}

func (d *Degree) FindDegree(src string, target string) {
	var err error
	var srcObj, trgtObj person
	for j := 0; j < 20; j++ {
		// Get the person's data
		srcObj, err = getPersonData(src)
		if err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if err != nil {
		fmt.Println("Invalid person url")
		return
	}
	for j := 0; j < 20; j++ {
		// Get the person's data
		trgtObj, err = getPersonData(target)
		if err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if err != nil {
		fmt.Println("Invalid person url")
		return
	}
	// Traversing through the person with less no of movies is efficient
	if len(trgtObj.Movies) > len(srcObj.Movies) {
		d.source = src
		d.target = target
	} else {
		d.source = target
		d.target = src
	}

	d.graphIn = make(chan result)
	d.out = make(chan bool)
	d.nextUrlsChan = make(chan []urls)
	d.urlsChan = make(chan []urls)

	d.personUrls = make(map[string]bool)
	d.movieUrls = make(map[string]bool)
	r := result{}
	r.currentUrls = append(r.currentUrls, urls{d.source, nil})
	go d.handleInput()
	d.graphIn <- r
	ok := <-d.out
	fmt.Println("Success =", ok)
	if ok {
		fmt.Println("Degree of seperation:", d.degree)
		fmt.Println("Connections :")
		for i := range d.connections {
			fmt.Println(i+1, ".Movie:", d.connections[i].movie)
			fmt.Println(d.connections[i].firstRole, ":", d.connections[i].first)
			fmt.Println(d.connections[i].secondRole, ":", d.connections[i].second)
		}
	}
}

func (d *Degree) isPersonParsed(url string) bool {
	d.Lock()
	res := d.personUrls[url]
	d.Unlock()
	return res
}

func (d *Degree) isMovieParsed(url string) bool {
	d.Lock()
	res := d.movieUrls[url]
	d.Unlock()
	return res
}

func (d *Degree) handleMovie(url string, pUrl string, pName string, pRole string, degree int, conns []connection) {
	var m movie
	var err error
	for j := 0; j < 20; j++ {
		// Get the movie data
		m, err = getMovieData(url)
		if err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if err != nil {
		var newUrls []urls
		d.nextUrlsChan <- newUrls
	} else {
		var newUrls []urls
		// Iterate through the list of casts
		for cst := range m.Cast {
			if d.isPersonParsed(m.Cast[cst].Url) != true {

				d.Lock()
				mainDegree := d.degree
				d.Unlock()
				// Sanity check(If the currentdegree is greater than degree no need to continue)
				if degree < mainDegree || mainDegree == 0 {
					// Leave the current person, might lead to infinite loop if not checked
					if m.Cast[cst].Url != pUrl {
						// If the current cast is the target
						if m.Cast[cst].Url == d.target {
							var connections []connection
							connections = append(connections, conns...)
							connections = append(connections, connection{m.Name, pName, pRole, m.Cast[cst].Name, m.Cast[cst].Role})
							d.Lock()
							d.degree = degree
							d.connections = connections
							d.Unlock()
							// close the channel
							d.out <- true
							close(d.graphIn)
							break
						} else {
							var connections []connection
							connections = append(connections, conns...)
							connections = append(connections, connection{m.Name, pName, pRole, m.Cast[cst].Name, m.Cast[cst].Role})

							u := urls{m.Cast[cst].Url, connections}
							newUrls = append(newUrls, u)
						}
					}
				}
			}
		}
		// Iterate through the list of crews
		for crw := range m.Crew {
			if d.isPersonParsed(m.Crew[crw].Url) != true {
				d.Lock()
				mainDegree := d.degree
				d.Unlock()
				// Sanity check(If the currentdegree is greater than degree no need to continue)
				if degree < mainDegree || mainDegree == 0 {
					// Leave the current person, might lead to infinite loop if not checked
					if m.Crew[crw].Url != pUrl {
						// If the current crew is the target
						if m.Crew[crw].Url == d.target {
							var connections []connection
							connections = append(connections, conns...)
							connections = append(connections, connection{m.Name, pName, pRole, m.Crew[crw].Name, m.Crew[crw].Role})
							d.Lock()
							d.degree = degree
							d.connections = connections
							d.Unlock()
							// close the channel
							d.out <- true
							close(d.graphIn)
							break
						} else {
							var connections []connection
							connections = append(connections, conns...)
							connections = append(connections, connection{m.Name, pName, pRole, m.Crew[crw].Name, m.Crew[crw].Role})
							u := urls{m.Crew[crw].Url, connections}
							newUrls = append(newUrls, u)

						}
					}
				}
			}
		}
		d.nextUrlsChan <- newUrls
	}
}

func (d *Degree) handlePerson(url string, degree int, connections []connection) {
	defer d.wg.Done()
	var err error
	var per person
	var nextUrls []urls

	for j := 0; j < 20; j++ {
		// Get the person's data
		per, err = getPersonData(url)
		if err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if err != nil {

	} else {
		// Iterate through all his movies
		mvCount := 0
		for mv := range per.Movies {
			// Check if the movie is already traversed
			if d.isMovieParsed(per.Movies[mv].Url) != true {
				d.Lock()
				d.movieUrls[per.Movies[mv].Url] = true
				d.Unlock()
				mvCount++
				go d.handleMovie(per.Movies[mv].Url, url, per.Name, per.Movies[mv].Role, degree, connections)
			}
		}

		for {
			if mvCount == 0 {
				break
			}
			select {
			case tempUrls, ok := <-d.nextUrlsChan:
				if !ok {
					break
				} else {
					mvCount--
					if len(tempUrls) > 0 {
						nextUrls = append(nextUrls, tempUrls...)
					}
				}
			default:
				time.Sleep(1 * time.Microsecond)
			}
		}

	}
	d.urlsChan <- nextUrls
}

func (d *Degree) findDegree(r result) {
	r.currentDegree++
	noPerson := 0
	// Iterate through all the users
	for i := range r.currentUrls {
		// Check if the person is already traversed
		if d.isPersonParsed(r.currentUrls[i].url) != true {
			noPerson++
			d.Lock()
			d.personUrls[r.currentUrls[i].url] = true
			d.Unlock()
			d.wg.Add(1)
			go d.handlePerson(r.currentUrls[i].url, r.currentDegree, r.currentUrls[i].connections)
		}
	}

	d.nextUrls = nil
	for {
		if noPerson == 0 {
			break
		}
		select {
		case tempUrls, ok := <-d.urlsChan:
			if !ok {
				break
			} else {
				noPerson--
				if len(tempUrls) > 0 {
					d.Lock()
					d.nextUrls = append(d.nextUrls, tempUrls...)
					d.Unlock()
				}
			}
		default:
			time.Sleep(1 * time.Microsecond)
		}
	}
	d.wg.Wait()

	if len(d.nextUrls) > 0 {
		// Send the url list again(Recursion)
		res := result{d.nextUrls, r.currentDegree}
		d.nextUrls = nil
		d.graphIn <- res
	} else {
		d.out <- false
	}
}

// Recieves url list and processes it
func (d *Degree) handleInput() {
	for {
		// Get the url list
		select {
		case r, ok := <-d.graphIn:
			// Check if the channel is closed
			if !ok {
				d.out <- false
				break
			} else {
				go d.findDegree(r)
			}
		default:
			time.Sleep(1 * time.Microsecond)
		}
	}
}

func parseConfig() error {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = json.Unmarshal(data, &conf)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func main() {
	err := parseConfig()
	if err == nil {
		// set the amount of CPU to be used
		if conf.NumCPU > 0 {
			runtime.GOMAXPROCS(conf.NumCPU)
		} else {
			numCPU := runtime.NumCPU()
			runtime.GOMAXPROCS((numCPU * 3) / 4)
		}
		// sanity check
		if len(os.Args) >= 3 {
			// No point calculating degree between the same person
			if os.Args[1] != os.Args[2] {
				d := new(Degree)
				start := time.Now()
				d.FindDegree(os.Args[1], os.Args[2])
				fmt.Println("Time taken to find degree =", time.Since(start).Seconds(), " Secs")
			} else {
				fmt.Println("Both are the same person, please give different person urls")
			}
		} else {
			fmt.Println("Please provide sufficient args")
		}
	}
}
