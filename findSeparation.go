/*
** Chandhana Surabhi 26/11/2019
 */

package main

import (
	"encoding/json"
	jsonb "encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/mtmoses/httprouter"
)

/*
GeneralInfo is the core structure for general information of movies or people
*/
type GeneralInfo struct {
	URL    string  `json:"url"`
	Type   string  `json:"type"`
	Name   string  `json:"name"`
	Movies []Movie `json:"movies"`
	Cast   []Cast  `json:"cast"`
	Crew   []Crew  `json:"crew"`
}

/*
Movie is the core structure for movie information
*/
type Movie struct {
	MovieName string `json:"name"`
	MovieURL  string `json:"url"`
	Role      string `json:"role"`
}

/*
Cast is the core struct for casts of the movie
*/
type Cast struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Role string `json:"role"`
}

/*
Crew is the core object for crew of the movie
*/
type Crew struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Role string `json:"role"`
}

/*
Server is the core structure for http router
*/
type Server struct {
	router *httprouter.Router
}

/*
Request is the core structure for web service input
*/
type Request struct {
	InputOne string `json:"inputone"`
	InputTwo string `json:"inputtwo"`
}

/*
Response is the core structure for web service output
*/
type Response struct {
	Status  bool     `json:"status"`
	Data    []Result `json:"data"`
	Degree  int      `json:"percentage"`
	Message string   `json:"message"`
}

/*
Result is the core object for output
*/
type Result struct {
	Movie           string `json:"movie"`
	Actor           string `json:"actor,omitempty"`
	SupportingActor string `json:"supportingactor,omitempty"`
	Role            string `json:"role"`
}

/*
global variables
*/
var degree int
var result []Result
var cast, crew []string

/*
SetveHTTP - setup http connection for webservices
*/
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Access-Control-Allow-Origin, Token, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Allow-Headers, *")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	s.router.ServeHTTP(w, r)
}

/*
showSplashscreen - test function to test API health
*/
func showSplashscreen() {
	screenImage := `
	API HEALTHY
`
	fmt.Println(screenImage)
	fmt.Println("===============")
	fmt.Println("API")
	fmt.Println("===============")
}

/*
healthCheckHandler - test function to test API health
*/

func healthCheckHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintln(w, "HI the api is working")
}

func initializeRoutes() {
	port := "8050"
	url := "localhost"

	portString := ":" + port
	fmt.Println("Starting server on\n", url, portString)

	router := httprouter.New()
	router.GET("/", healthCheckHandler)
	router.POST("/user/v1/checkdegree", checkDegreeHandler)

	http.ListenAndServe(":8050", &Server{router})
}

/*
checkDegreeHandler - checks the degree of separation
*/
func checkDegreeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var names *Request

	//Reading request from the body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Parsing Error", http.StatusInternalServerError)
	}

	//Unmarshaling the request to movies struct
	err = json.Unmarshal(body, &names)
	if err != nil {
		http.Error(w, "Parsing Error", http.StatusInternalServerError)
		return
	}

	degree = 0

	/*
		get first degree relation - directly involved in movies
	*/
	movieData := getData(names.InputOne)
	getFirstDegreeRelation(movieData, names.InputOne, names.InputTwo)

	if degree <= 1 {
		goto resultResponse
	}

	/*
		get second degree relation - same cast and crew with both people
	*/
	movieData = getData(names.InputTwo)
	getSecondDegreeRelation(movieData, names.InputTwo, names.InputOne)

	if degree <= 2 {
		goto resultResponse
	}
	/*
		get third degree relation - related with some people in between
	*/
	getThirdDegreeRelation(names.InputOne, names.InputTwo)
	fmt.Println(degree)

	/*
		make json response structure
	*/
resultResponse:
	res2D := Response{
		Message: "degree of separation collected",
		Status:  true,
		Data:    result,
		Degree:  degree,
	}
	fmt.Fprintln(w, jSONResponse(res2D))
	return
}

func main() {
	showSplashscreen()
	initializeRoutes()
}

/*
jsonResponse constructs the json response
*/
func jSONResponse(resp Response) string {
	j, err := json.Marshal(&resp)
	if err != nil {
	}
	return string(j)
}

/*
getData gets the movie info from moviebuff
*/
func getData(name string) (info GeneralInfo) {
	movieURL := fmt.Sprintf("https://data.moviebuff.com/%s", strings.Replace(strings.ToLower(name), " ", "-", -1))
	fmt.Println(movieURL)
	req, _ := http.NewRequest("GET", movieURL, nil)
	resp, _ := http.DefaultClient.Do(req)
	decoder := jsonb.NewDecoder(resp.Body)
	decoder.Decode(&info)

	return info
}

/*
getFirstDegreeRelation gets the relation between people - directly involved in a movies
*/
func getFirstDegreeRelation(movies GeneralInfo, personOne, personTwo string) {

	for i, _ := range movies.Movies {
		movieInfo := getData(movies.Movies[i].MovieURL)
		getInfoOfMovie(movieInfo, personOne, personTwo)
		if degree >= 1 {
			return
		}
	}
}

/*
getSecondDegreeRelation gets the relation between people - cast and crew worked with both the peoples searched for
*/
func getSecondDegreeRelation(movies GeneralInfo, personOne, personTwo string) {

	for i, _ := range movies.Movies {
		movieInfo := getData(movies.Movies[i].MovieURL)
		for j, _ := range cast {
			getInfoOfMovie(movieInfo, "", cast[j])
			if degree >= 2 {
				return
			}
		}
		for k, _ := range crew {
			getInfoOfMovie(movieInfo, "", crew[k])
			if degree >= 2 {
				return
			}
		}
	}
}

/*
getInfoOfMovie - find info of movie and compare the cast and crew
*/
func getInfoOfMovie(movie GeneralInfo, personOne, personTwo string) {
	var resultSet Result

	for j, _ := range movie.Cast {

		if strings.ToLower(movie.Cast[j].Name) == strings.ToLower(personTwo) {
			degree = degree + 1
			resultSet.Actor = movie.Cast[j].Name
			resultSet.Movie = movie.Name
			resultSet.Role = movie.Cast[j].Role
			result = append(result, resultSet)
			return
		} else {
			if !contains(cast, strings.ToLower(movie.Cast[j].Name)) {
				cast = append(cast, strings.ToLower(movie.Cast[j].Name))
			}
		}
	}
	for k, _ := range movie.Crew {
		if strings.ToLower(movie.Crew[k].Name) == strings.ToLower(personTwo) {
			degree = degree + 1
			resultSet.Actor = movie.Crew[k].Name
			resultSet.Movie = movie.Name
			resultSet.Role = movie.Crew[k].Role
			result = append(result, resultSet)
			return
		} else {
			if !contains(crew, strings.ToLower(movie.Crew[k].Name)) {
				crew = append(crew, strings.ToLower(movie.Crew[k].Name))
			}
		}
	}

}

/*
getThirdDegreeRelation gets the relation between two people not connected directly
*/
func getThirdDegreeRelation(personOne, personTwo string) {

	for i, _ := range cast {

		movieItem := getData(cast[i])
		getFirstDegreeRelation(movieItem, "", cast[i])
		getSecondDegreeRelation(movieItem, "", personTwo)

		if degree > 2 {
			degree = 3
			return
		}
	}
	for j, _ := range crew {
		movieItem := getData(crew[j])
		getFirstDegreeRelation(movieItem, "", crew[j])
		getSecondDegreeRelation(movieItem, "", personTwo)
		if degree > 2 {
			degree = 3
			return
		}
	}
}

/*contains checks whether the given string is already present in the slice
Gets string slice as input
Returns bool
*/
func contains(slice []string, element string) bool {
	for _, i := range slice {
		if i == element {
			return true
		}
	}
	return false
}
