package main

import "net/http"

func unpad(buf []byte) []byte {
	// Assume valid length and padding. Should add checks
	padding := int(buf[len(buf)-1])
	return buf[:len(buf)-padding]
}

func main() {
	fs := http.FileServer(http.Dir("E:/phpstudy_pro/WWW"))
	http.Handle("/", fs)
	http.ListenAndServe(":8999", nil)
}
