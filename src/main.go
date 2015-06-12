package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var out io.Writer = ioutil.Discard
var FileLogger *log.Logger

func init() {
	startLogger()
}

func main() {
	FileLogger.Println("start time::", time.Now())
	FileLogger.Println("inside main")
	args := os.Args[1:]
	FileLogger.Println("got arguments::", args)
	GetSeparation(args[0], args[1])
	FileLogger.Println("its Done!!!")
	FileLogger.Println("end time::", time.Now())
}

func startLogger() {
	logMode := "file"
	var err error
	switch logMode {
	case "file":
		out, err = os.Create("/users/logs/degrees.log")
		if err != nil {
			fmt.Println(err)
		}
	case "screen":
		out = os.Stdout

	default:
		out = ioutil.Discard
	}
	FileLogger = log.New(out, "", log.Lshortfile)
}
