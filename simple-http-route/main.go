package main

import (
	"fmt"
	"net/http"
)

func main() {
	var r route
	http.ListenAndServe(":8080", &r)
}

type route struct{}

func (ro *route) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/a":
		fmt.Fprint(w, "Executing a")
	case "/b":
		fmt.Fprint(w, "Executing b")
	case "/c":
		fmt.Fprint(w, "Executing c")
	default:
		http.Error(w, "404 Not Found", 404)
	}
}
