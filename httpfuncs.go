package main

import (
	"fmt"
	"os"
)

const moviebuff = "http://data.moviebuff.com/"

func ErrHandle(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func main(){
	url := moviebuff + os.Args[1]

	json, err := getData(url)
	defer ErrHandle(err)

	fmt.Println(json.Name, json.Url)

}

