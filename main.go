package main

import (

	//"log"
	"log"
	"myproject/challenge2015/internal/handlers"

	//"net/http"

	"net/http"
)

func main() {

	err := http.ListenAndServe(":8080", handlers.New())
	if err != nil {
		log.Fatal("error ", err.Error())
	}
}
