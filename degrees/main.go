package main

import (
	"degrees/degreecalculator"
	"degrees/moviebuffclient"
	"fmt"
	"os"
)

func main() {
	// take input as commnad line arguments
	input := os.Args
	if len(input) < 3 {
		fmt.Println("Please provide two person urls as command line arguments")
		return
	}
	person1Url := input[1]
	person2Url := input[2]

	// validate the input
	dir, err := validateInputAndProvideFlowDirection(person1Url, person2Url)
	if err != nil {
		fmt.Println(fmt.Errorf("error while validating input: %s", err))
		return
	}

	// calculate the degree of separation
	if dir {
		separation, chainInfo, err := degreecalculator.Calculate(person1Url, person2Url)
		degreecalculator.PrintRespose(separation, chainInfo, err)
		return
	}
	separation, chainInfo, err := degreecalculator.Calculate(person2Url, person1Url)
	degreecalculator.PrintResponseInReverse(separation, chainInfo, err)

}

func validateInputAndProvideFlowDirection(person1Url string, person2Url string) (bool, error) {
	if person1Url == person2Url {
		return false, fmt.Errorf("given urls are same. Please provide two different urls")
	}
	p1Info, err := moviebuffclient.GetPersonInfo(person1Url)
	if err != nil {
		return false, err
	}
	p2Info, err := moviebuffclient.GetPersonInfo(person2Url)
	if err != nil {
		return false, err
	}
	if len(p1Info.Movies) < len(p2Info.Movies) {
		return true, nil
	}
	return false, nil
}
