package main

import (
	"code.google.com/p/go.net/websocket"
	"flag"
	"log"
	"net/http"
	"text/template"
)

var addr = flag.String("addr", ":12560", "ZoneServer Service")

func homeHandler(c http.ResponseWriter, req *http.Request) {

	path := "home.html"

	if len(req.URL.String()) > 1 {
		path = req.URL.String()[1:]
	}

	var templ = template.Must(template.ParseFiles(path))
	templ.Execute(c, req.Host)
}

func main() {
	flag.Parse()

	go h.run()

	http.Handle("/ws", websocket.Handler(wsHandler))

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
