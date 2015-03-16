package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
	//"github.com/trustmaster/goflow"
)

type Degree struct {
	source     string
	target     string
	degree     int
	res        result
	graphIn    chan result
	out        chan bool
	personUrls map[string]bool
	movieUrls  map[string]bool
	wg         sync.WaitGroup
	nextUrls   []string
}

func (d *Degree) FindDegree(src string, target string) {
	d.graphIn = make(chan result)
	d.out = make(chan bool)
	d.source = src
	d.target = target
	d.personUrls = make(map[string]bool)
	d.movieUrls = make(map[string]bool)
	//degreeGraph := CreateDegreeGraph(d, graphIn)

	//degreeGraph.SetInPort("In", graphIn)
	//flow.RunNet(degreeGraph)

	////close(graphIn)
	//<-degreeGraph.Wait()
	//fmt.Println("degree :", d.degree)
	//fmt.Println("Connections :", d.res)
	r := result{}
	r.currentUrl = append(r.currentUrl, src)
	go d.handleInput()
	d.graphIn <- r
	ok := <-d.out
	fmt.Println("Success =", ok)
	fmt.Println("Degree of seperation:", d.degree)
}

func (d *Degree) handlePerson(r *result, i int) {
	defer d.wg.Done()
	//	if r.currentUrl[i] == d.target {
	//		d.degree = r.currentDegree
	//		//		d.res = r
	//		d.out <- true
	//		return
	//	}
	var err error
	var per person
	for j := 0; j < 20; j++ {
		per, err = getPersonData(r.currentUrl[i])
		if err == nil {
			break
		}
	}
	if err != nil {
		//				d.out <- false
		//				return
	} else {
		for mv := range per.Movies {
			if d.movieUrls[per.Movies[mv].Url] != true {
				d.movieUrls[per.Movies[mv].Url] = true
				var m movie
				for j := 0; j < 20; j++ {
					m, err = getMovieData(per.Movies[mv].Url)
					if err == nil {
						break
					}
				}
				if err != nil {
					//							d.out <- false
					//							break
				} else {
					for cst := range m.Cast {
						if r.currentDegree < d.degree || d.degree == 0 {

							if m.Cast[cst].Url != r.currentUrl[i] {
								if m.Cast[cst].Url == d.target {
									d.degree = r.currentDegree
									fmt.Println(m.Cast[cst].Url)
									fmt.Println(m.Url)
									d.out <- true
									close(d.graphIn)
									break
								} else {
									//									fmt.Println(m.Cast[cst].Url)
									d.nextUrls = append(d.nextUrls, m.Cast[cst].Url)
								}
							}
						}
					}
					for crw := range m.Crew {
						if r.currentDegree < d.degree || d.degree == 0 {
							if m.Crew[crw].Url != r.currentUrl[i] {
								if m.Crew[crw].Url == d.target {
									d.degree = r.currentDegree
									fmt.Println(m.Crew[crw].Url)
									fmt.Println(m.Url)
									d.out <- true
									close(d.graphIn)
									break
								} else {
									d.nextUrls = append(d.nextUrls, m.Crew[crw].Url)
								}
							}
						}
					}
				}
			}
		}
	}
}

func (d *Degree) findDegree(r result) {
	fmt.Println(r.currentUrl)
	r.currentDegree++
	for i := range r.currentUrl {
		if d.personUrls[r.currentUrl[i]] != true {
			d.personUrls[r.currentUrl[i]] = true
			d.wg.Add(1)
			go d.handlePerson(&r, i)
		}
	}
	d.wg.Wait()
	if len(d.nextUrls) > 0 {
		res := result{d.nextUrls, r.currentDegree}
		d.nextUrls = nil
		d.graphIn <- res
	} else {
		d.out <- false
	}

}

func (d *Degree) handleInput() {
	for {
		select {
		case r, ok := <-d.graphIn:
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

func main() {
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS((numCPU * 3) / 4)

	if len(os.Args) >= 3 {
		if os.Args[1] != os.Args[2] {
			d := new(Degree)
			d.FindDegree(os.Args[1], os.Args[2])
		}
	} else {
		fmt.Println("Please provide sufficient args")
	}
}
