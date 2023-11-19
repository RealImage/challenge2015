package main

import (
	services "challenge2015/internal/domain/services"
	"fmt"
	"os"
)

func main() {
	args := os.Args

	if len(args) != 3 {
		fmt.Println("provide atleast two arguments")
		return
	}
	if args[1] == args[2] {
		fmt.Println("Please provide different actors")
		return
	}
	services.FirstActor = args[1]
	services.SecondActor = args[2]

	services.ActorListForFirst = append(services.ActorListForFirst, args[1])
	services.ActorListForSecond = append(services.ActorListForSecond, args[2])

	err := services.SmallestDegreeOfSeparation()
	if err != nil {
		fmt.Println("Error in finding smallest degree of separation: ", err)
	}
}
