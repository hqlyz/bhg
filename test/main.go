package main

import (
	"fmt"
)

func unpad(buf []byte) []byte {
	// Assume valid length and padding. Should add checks
	padding := int(buf[len(buf)-1])
	return buf[:len(buf)-padding]
}

func main() {
	// fs := http.FileServer(http.Dir("g:/books"))
	// http.Handle("/", fs)
	// http.ListenAndServe(":8080", nil)

	b := []byte{65, 66, 67, 68, 69}
	b2 := unpad(b)
	fmt.Println(b2)
}
