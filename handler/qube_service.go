package handler

import (
	"fmt"
	"strings"

	srvError "qubecinemas.com/error"
	utility "qubecinemas.com/internal/utility"
)

//Processing of data to find Minimum degree of Seperation between two people
/*
FLOW: Possibility
1. Actor -> Movie -> Second Person (if found) - Minimum Steps
2. Actor -> Movie -> Actor
3. Actor -> Movie -> Actor -> Movie -> Actor
4. *****

NOTE: Actor to Movie is ONE Single Step
*/

func GetProcessedData(firstPerson, secondPerson string, count int) (int, error) {
	secondPerson = strings.Join(strings.Split(secondPerson, "-"), " ")
	if strings.EqualFold(strings.ToLower(firstPerson), strings.ToLower(secondPerson)) {
		fmt.Println("---------------------------------------------------------------")
		return count, nil
	}

	//First Person to get movie names
	err, person := utility.PersonUtility(firstPerson)
	if err != nil {
		return 0, fmt.Errorf(srvError.PackError(&srvError.UnableToBindResourceBody, err))
	}
	for _, value := range person.Movies {
		err, movieSub := utility.MovieUtility(value.Url)
		if err != nil {
			return 0, fmt.Errorf(srvError.PackError(&srvError.UnableToBindResourceBody, err))
		}
		people := append(movieSub.Cast, movieSub.Crew...)
		count += 1
		for _, inval := range people {
			fmt.Println(strings.ToLower(inval.Name), secondPerson)
			GetProcessedData(strings.ToLower(inval.Name), secondPerson, count)
		}

	}
	return count, nil
}
