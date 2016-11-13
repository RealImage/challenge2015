package main 

import (
	// "encoding/json"
	// "fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	// "os"
	// "strconv"
	// "time"
	// "strings"
)

func main() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", HandleIndex)
	// rtr.HandleFunc("/{user:[a-zA-Z0-9]+}/{pass:[a-zA-Z0-9]+}/GetSequencevalue", get_sequencevalue).Methods("GET")

	http.Handle("/", rtr)
	log.Println(http.ListenAndServe(":3011", nil))
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	const tpl = `
		<!DOCTYPE html>
		<html>
			<head>
				<meta charset="UTF-8">
				<title>Degrees Of Separation</title>
			</head>
			<body>
				<h2>Degrees Of Separation</h2>
				<form action="">
				  First name:<br>
				  <input type="text" name="firstname" value="vijay">
				  <br>
				  Last name:<br>
				  <input type="text" name="lastname" value="ajith-kumar">
				  <br><br>
				  <input type="submit" value="Submit">
				  <input type="button" value="Check!" onclick="checkDoS()">
				</form>
				<script type="text/javascript">
				  	function checkDoS(){
					    alert("Degrees Of Separation");
					}
				</script>
			</body>
		</html>`
	/*check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}*/
	io.WriteString(w, tpl)
}