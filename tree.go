package main

type dataset struct{
	initialUrl string
	finalUrl string
	root	*dataset
	branch  *dataset
	brMov map[string]movie
}
