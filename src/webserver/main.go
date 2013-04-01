package main

import (
	"flag"
	"log"
	"net/http"
	"io/ioutil"
)

var addr = flag.String("addr", ":80", "WebServer Service")

func pageHandler(c http.ResponseWriter, req *http.Request) {

	path := "home.html"

	if len(req.URL.String()) > 1 {
		path = req.URL.String()[1:]
	}

	b, err := ioutil.ReadFile(path)

	if err != nil {
		return
	}

	c.Write(b)
}

func main() {
	flag.Parse()

	http.HandleFunc("/", pageHandler)

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}