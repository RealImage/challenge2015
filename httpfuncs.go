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
	if len(os.Args) != 3 {
		fmt.Print("Usage Example : degrees vn-mayekar magie-mathur")
	}
	buildTree(os.Args[1], os.Args[2], os.Args[1], true)
}

func buildTree(argument, destination, parent string, getMovie bool ){

	url := moviebuff + argument
	json, err := getData(url)
	defer ErrHandle(err)

	if getMovie {
		parent = argument
	}

	for _, movie := range json.Movies{
		fmt.Println(movie.Url)
		buildTree(movie.Url, destination, argument, false)
	}

	for _, cast := range json.Cast{
		fmt.Println(cast.Url)
	}

}
