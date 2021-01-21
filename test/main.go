package main

import (
	"encoding/hex"
	"fmt"
)

func main() {
	// fs := http.FileServer(http.Dir("g:/books"))
	// http.Handle("/", fs)
	// http.ListenAndServe(":8080", nil)

	r, _ := hex.DecodeString("abcd中")
	fmt.Println(r)
}
