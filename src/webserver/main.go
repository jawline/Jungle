package main

import (
	"flag"
	"log"
	"net/http"
	"io/ioutil"
	"fmt"
)

var addr = flag.String("addr", ":80", "WebServer Service")
var fileDirectory = ""

func pageHandler(c http.ResponseWriter, req *http.Request) {

	path := "home.html"

	if len(req.URL.String()) > 1 {
		path = req.URL.String()[1:]
	}

	b, err := ioutil.ReadFile(fileDirectory + path)

	if err != nil {
		fmt.Printf("Error: " + err.Error())
		return
	}

	c.Write(b)
}

func main() {

	directoryFlag := flag.String("path", "./", "The path to the server files")
	flag.Parse()
	fileDirectory = *directoryFlag

	http.HandleFunc("/", pageHandler)

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
