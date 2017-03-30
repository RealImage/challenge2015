package main 

import(
	"os"
	"fmt"
	"log"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type movies struct {
	Name string 		`json:"name"`
	Url string			`json:"url"`
	Role string 		`json:"role"`
}

type peopleInfo struct {
	Url string			`json:"url"`
	Type string			`json:"type"`
	Name string			`json:"name"`
	Movies []movies `json:"movies"`
	Cast   []movies  `json:"cast"`
	Crew   []movies  `json:"crew"`
}

type dos struct {
	movie string
	people1 string
	role1 string
	people2 string
	role2 string
}

type MovieBuffs struct {
	source string
	destination string
	people1 peopleInfo
	people2 peopleInfo
	link map[string]dos
	visit []string
	visited map[string]bool
	p2Movies map[string]movies
}
var movieBuff MovieBuffs
func getData(url string) (peopleInfo, error) {
	var pInfo peopleInfo
	resp, _ := http.Get("http://data.moviebuff.com/" + url)
	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return pInfo, err
	}

	err = json.Unmarshal(res, &pInfo)
	if err != nil {
		return pInfo, err
	}

	return pInfo, nil
}

func getPeopleInfo(people1, people2 string) error {
	
	movieBuff.p2Movies, movieBuff.visited, movieBuff.link = make(map[string]movies), make(map[string]bool), make(map[string]dos)
	p1, err := getData(people1)
	if err != nil {
		return err
	}
	p2, err := getData(people2)
	if err != nil {
		return err
	}

	if len(p1.Movies) > len(p2.Movies) {
		movieBuff.source, movieBuff.destination = people2, people1
		movieBuff.people1, movieBuff.people2 = p2, p1
	} else {
		movieBuff.source, movieBuff.destination = people1, people2
		movieBuff.people1, movieBuff.people2 = p1, p2
	}

	for _, movie := range movieBuff.people2.Movies {
		movieBuff.p2Movies[movie.Url] = movie
	}
	movieBuff.visit = append(movieBuff.visit, movieBuff.source)
	movieBuff.visited[movieBuff.source] = true
	dos, err := findDos()
	if err != nil {
		log.Fatalln(err.Error())
	}

	//result
	fmt.Printf("Degree of separation: %d \n", len(dos))
	for i, d := range dos {
		fmt.Printf("\n%d. ",i+1)
		fmt.Println("Movie: "+ d.movie)
		fmt.Println(d.role1+": "+ d.people1)
		fmt.Println(d.role2+": "+ d.people2)
	}

	return nil
}

func findDos() ([]dos, error) {

	var d []dos
	for true {

		for _, person := range movieBuff.visit {

			people1, err := getData(person)
			if err != nil {
				if strings.Contains(err.Error(), "looking for beginning of value") {
					continue
				}
				return nil, err
			}

			for _, p1movie := range people1.Movies {
				if movieBuff.p2Movies[p1movie.Url].Url == p1movie.Url {
					if _, found := movieBuff.link[people1.Url]; found {
						d = append(d, movieBuff.link[people1.Url], dos{p1movie.Name, people1.Name, p1movie.Role, movieBuff.people2.Name, movieBuff.p2Movies[p1movie.Url].Role})
					} else {
						d = append(d, dos{p1movie.Name, people1.Name, p1movie.Role, movieBuff.people2.Name, movieBuff.p2Movies[p1movie.Url].Role})
					}
					return d, nil
				}
			}

			// Find new nodes to continue searching
			for _, p1movie := range people1.Movies {

				if movieBuff.visited[p1movie.Url] {
					continue
				}

				movieBuff.visited[p1movie.Url] = true

				p1moviedetail, err := getData(p1movie.Url)
				if err != nil {
					if strings.Contains(err.Error(), "looking for beginning of value") {
						continue
					}
					return nil, err
				}
				for _, p1moviecast := range p1moviedetail.Cast {

					if movieBuff.visited[p1moviecast.Url] {
						continue
					}

					movieBuff.visited[p1moviecast.Url] = true
					movieBuff.visit = append(movieBuff.visit, p1moviecast.Url)
					movieBuff.link[p1moviecast.Url] = dos{p1movie.Name, people1.Name, p1movie.Role, p1moviecast.Name, p1moviecast.Role}
				}

				for _, p1moviecrew := range p1moviedetail.Crew {

					if movieBuff.visited[p1moviecrew.Url] {
						continue
					}

					movieBuff.visited[p1moviecrew.Url] = true
					movieBuff.visit = append(movieBuff.visit, p1moviecrew.Url)
					movieBuff.link[p1moviecrew.Url] = dos{p1movie.Name, people1.Name, p1movie.Role, p1moviecrew.Name, p1moviecrew.Role}
				}

			}
		}

	}

	return d, nil
}

func main() {
	args := os.Args[1:]

	if len(args) != 2 {
		log.Fatalln("Please pass the two arguments")
	}
	err := getPeopleInfo(args[0], args[1])
	if err != nil {
		log.Fatalln(err)
	}
}
