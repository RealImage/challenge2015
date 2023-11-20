package main

import (
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

var INVALID_VAL int = 100005
var LOGS_BATCH_SIZE = 1000

func main() {

	args := os.Args

	if len(args) < 3 {
		fmt.Println("Invalid inputs. Input format eg. \"amitabh-bachchan robert-de-niro\"")
		return
	}

	app := GlobalVar{
		IsNodeVisitedSet: make(map[string]bool),
		FirstPersonNode:  InfoNodeForQueue{InfoNodeEntry: InfoNode{URL: args[1]}},
		SecondPersonNode: InfoNodeForQueue{InfoNodeEntry: InfoNode{URL: args[2]}},
	}
	app.findDegreesOfSeparation()

	app.PrintAns()
}

type GlobalVar struct {
	IsNodeVisitedSet map[string]bool
	FirstPersonNode  InfoNodeForQueue
	SecondPersonNode InfoNodeForQueue
}

func (app *GlobalVar) findDegreesOfSeparation() error {
	sFunctionName := "findDegreesOfSeparation"

	if len(app.FirstPersonNode.InfoNodeEntry.URL) == 0 || len(app.SecondPersonNode.InfoNodeEntry.URL) == 0 {
		fmt.Println(sFunctionName, "Please provide valid inputs")
		return errors.New("Invalid data provided")
	}

	infoNodeQueue := list.New()

	app.addNodeToQueue(infoNodeQueue, InfoNodeForQueue{
		ParentNodeEntry: &InfoNodeForQueue{},
		InfoNodeEntry: InfoNode{
			URL: app.FirstPersonNode.InfoNodeEntry.URL,
		}})

	bIsCurrentPerson := true

	for false == app.checkAnsFound() && infoNodeQueue.Len() > 0 {

		nextInfoNodeQueue := list.New()

		for infoNodeQueue.Len() > 0 && false == app.checkAnsFound() {
			removedNode := infoNodeQueue.Front()
			infoNodeQueue.Remove(removedNode)

			if nil == removedNode.Value {
				fmt.Println(sFunctionName, "Found invalid value during traversing at level")
				return errors.New(fmt.Sprintf("Found invalid value during traversing at level"))
			}

			castedRemovedNode := removedNode.Value.(InfoNodeForQueue)

			if len(castedRemovedNode.InfoNodeEntry.URL) == 0 {
				fmt.Println(sFunctionName, "Found invalid value during traversing at level")
				return errors.New(fmt.Sprintf("Found invalid value during traversing at level"))
			}

			app.populateNeighbours(castedRemovedNode, nextInfoNodeQueue, bIsCurrentPerson)

			if nextInfoNodeQueue.Len()%LOGS_BATCH_SIZE == 0 {
				fmt.Printf("[%d] items added to next queue\n", nextInfoNodeQueue.Len())
			}
		}

		if false == app.checkAnsFound() {
			// current BFS done, use the new infoNodeQueue for next iteration
			infoNodeQueue = nextInfoNodeQueue

			// next time we will be iterating opposite
			bIsCurrentPerson = !bIsCurrentPerson
		}
	}

	return nil
}

func (app *GlobalVar) populateNeighbours(infoNode InfoNodeForQueue, queue *list.List, bIsCurrentPerson bool) error {
	sFunctionName := "populateNeighbours"

	if len(infoNode.InfoNodeEntry.URL) == 0 {
		fmt.Println(sFunctionName, "Invalid Moviebuff URL provided")
		return errors.New("Invalid Moviebuff URL provided")
	}

	sFormattedURL := fmt.Sprintf("http://data.moviebuff.com/%s", infoNode.InfoNodeEntry.URL)

	// Make GET request
	response, err := app.doHTTPRequest(sFormattedURL)

	if err != nil {
		fmt.Println(sFunctionName, "Error making the request:", err)
		return err
	}

	if bIsCurrentPerson {
		app.unmarshalAndPopulatePerson(response, queue, infoNode)
	} else {
		app.unmarshalAndPopulateMovies(response, queue, infoNode)
	}

	return nil
}

func (app *GlobalVar) unmarshalAndPopulatePerson(rawData []byte, queue *list.List, parentInfoNode InfoNodeForQueue) error {
	// sFunctionName := "unmarshalAndPopulatePerson"

	// Unmarshal JSON into struct
	var responseData MovieBuffResponseDataForPerson
	err := json.Unmarshal(rawData, &responseData)
	if err != nil {
		// fmt.Println(sFunctionName, "Error unmarshalling JSON:", err)
		return err
	}

	if parentInfoNode.InfoNodeEntry.URL == app.FirstPersonNode.InfoNodeEntry.URL {
		app.FirstPersonNode.InfoNodeEntry = InfoNode{
			URL:  responseData.URL,
			Name: responseData.Name,
		}

		parentInfoNode.InfoNodeEntry = InfoNode{
			URL:  responseData.URL,
			Name: responseData.Name,
		}
	}

	for _, node := range responseData.Movies {
		toAdd := InfoNodeForQueue{
			ParentNodeEntry: &parentInfoNode,
			InfoNodeEntry:   node,
		}
		app.addNodeToQueue(queue, toAdd)
	}

	return nil
}

func (app *GlobalVar) unmarshalAndPopulateMovies(rawData []byte, queue *list.List, parentInfoNode InfoNodeForQueue) error {
	// sFunctionName := "unmarshalAndPopulatePerson"

	// Unmarshal JSON into struct
	var responseData MovieBuffResponseDataForMovies
	err := json.Unmarshal(rawData, &responseData)
	if err != nil {
		// fmt.Println(sFunctionName, "Error unmarshalling JSON:", err)
		return err
	}

	for _, node := range responseData.Cast {
		toAdd := InfoNodeForQueue{
			ParentNodeEntry: &parentInfoNode,
			InfoNodeEntry:   node,
		}
		app.addNodeToQueue(queue, toAdd)
	}

	for _, node := range responseData.Crew {
		toAdd := InfoNodeForQueue{
			ParentNodeEntry: &parentInfoNode,
			InfoNodeEntry:   node,
		}
		app.addNodeToQueue(queue, toAdd)
	}

	return nil
}

func (app *GlobalVar) addNodeToQueue(queue *list.List, node InfoNodeForQueue) bool {

	if app.IsNodeVisitedSet[node.InfoNodeEntry.URL] {
		return false
	}

	app.IsNodeVisitedSet[node.InfoNodeEntry.URL] = true
	queue.PushBack(node)

	if node.InfoNodeEntry.URL == app.SecondPersonNode.InfoNodeEntry.URL {
		app.SecondPersonNode = node
	}

	if queue.Len()%LOGS_BATCH_SIZE == 0 {
		fmt.Printf("[%d] items added to queue\n", queue.Len())
	}

	return true
}

func (app *GlobalVar) checkAnsFound() bool {
	return app.IsNodeVisitedSet[app.SecondPersonNode.InfoNodeEntry.URL]
}

func (app *GlobalVar) PrintAns() {

	if app.SecondPersonNode.ParentNodeEntry == nil {
		fmt.Println("Unable to find the answer")
		return
	}

	type AnsNodes struct {
		Movie        string
		FirstPerson  string
		SecondPerson string
	}

	list := list.New()

	temp := app.SecondPersonNode

	for temp.ParentNodeEntry.InfoNodeEntry.URL != "" {
		curAnsNodes := AnsNodes{
			Movie:        temp.ParentNodeEntry.InfoNodeEntry.Name,
			SecondPerson: temp.InfoNodeEntry.Role + ": " + temp.InfoNodeEntry.Name,
			FirstPerson:  temp.ParentNodeEntry.InfoNodeEntry.Role + ": " + temp.ParentNodeEntry.ParentNodeEntry.InfoNodeEntry.Name,
		}

		temp = *(*temp.ParentNodeEntry).ParentNodeEntry

		list.PushFront(curAnsNodes)
	}

	index := 1

	fmt.Println("Degrees of Separation: ", list.Len())

	for element := list.Front(); element != nil; element = element.Next() {
		node := element.Value.(AnsNodes)

		fmt.Printf("%d. Movie: %s\n", index, node.Movie)
		fmt.Println(node.FirstPerson)
		fmt.Println(node.SecondPerson)

		fmt.Println()

		index++
	}
}
