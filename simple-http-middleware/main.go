package main

import (
	"fmt"
	"net/http")

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello %s\n", r.URL.Query().Get("name"))
}

func main() {
	f := http.HandlerFunc(hello)
	var l *logger = &logger{
		Inner: f,
	}
	http.ListenAndServe(":8080", l)
}

type logger struct {
	Inner http.Handler
}

func (l *logger)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "start")
	l.Inner.ServeHTTP(w, r)
	fmt.Fprintln(w, "end")
}