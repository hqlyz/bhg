package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func main() {
	for i := 1; i <= 3; i++ {
		fmt.Println("Try", i, "...")
		conn, err := net.Dial("tcp", "192.168.206.175:21")
		if err != nil {
			fmt.Printf("Error occurred at offset %d: %v\n", i, err)
		}
		buf := bufio.NewReader(conn)
		str, _ := buf.ReadString('\n')
		fmt.Println(str)

		fmt.Fprintf(conn, "USER %s\n", strings.Repeat("A", i))
		str, _ = buf.ReadString('\n')
		fmt.Println(str)

		fmt.Fprint(conn, "PASS password\n")
		str, _ = buf.ReadString('\n')
		fmt.Println(str)

		if err = conn.Close(); err != nil {
			fmt.Printf("Error occurred at offset %d: %v\n", i, err)
		}
		fmt.Println()
	}
}
