package main

import (
	"debug/pe"
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

func main() {
	// D:\\Program Files\\Telegram Desktop\\Telegram.exe
	f, err := os.Open("D:\\Program Files\\Telegram Desktop\\Telegram.exe")
	if err != nil {
		log.Fatalln(err)
	}
	pefile, err := pe.NewFile(f)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	defer pefile.Close()

	dosHeader := make([]byte, 96)
	sizeOffset := make([]byte, 4)
	_, err = f.Read(dosHeader)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("[-----DOS Header / Stub-----]")
	fmt.Printf("[+] Magic Value: %s%s\n", string(dosHeader[0]), string(dosHeader[1]))
	sigOffset := binary.LittleEndian.Uint32(dosHeader[0x3c:])

	_, err = f.ReadAt(sizeOffset, int64(sigOffset))
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("[-----Signature Header-----]")
	fmt.Printf("[+] LFANEW Value: %s\n", string(sizeOffset))
}
