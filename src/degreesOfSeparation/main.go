package main 

import (
	// "encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"degreesOfSeparation/html"
	// "os"
	// "strconv"
	// "time"
	"strings"
)

func main() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", HandleIndex)
	rtr.HandleFunc("/checkDoS", checkDoS).Methods("POST")	
	http.Handle("/", rtr)
	log.Println(http.ListenAndServe(":3011", nil))
}
/**
 * [HandleIndex description]
 * @param {[type]} w http.ResponseWriter [description]
 * @param {[type]} r *http.Request       [description]
 */
func HandleIndex(w http.ResponseWriter, r *http.Request) {	
	tpl := html.Template
	io.WriteString(w, tpl)
}
/**
 * [checkDoS description]
 * @param  {[type]} w http.ResponseWriter [description]
 * @param  {[type]} r *http.Request       [description]
 * @return {[JSON]}   [DoS]
 */
func checkDoS(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	if r.Method == "GET" {
		fmt.Fprintf(w, "Error Method call")
	} else {
		actor1 := strings.TrimSpace(r.PostFormValue("actor1"))
		actor2 := strings.TrimSpace(r.PostFormValue("actor2"))
		w.Write([]byte(actor1+" DoS "+actor2))
	}	
}