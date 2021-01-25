package main

import (
	"debug/pe"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	// D:\\Program Files\\Telegram Desktop\\Telegram.exe
	f, err := os.Open("E:\\Program Files\\Telegram Desktop\\Telegram.exe")
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
	sigOffset := int64(binary.LittleEndian.Uint32(dosHeader[0x3c:]))

	_, err = f.ReadAt(sizeOffset, sigOffset)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("[-----Signature Header-----]")
	fmt.Printf("[+] LFANEW Value: %s\n", string(sizeOffset))

	// Create the reader and read COFF Header
	sr := io.NewSectionReader(f, 0, 1<<63-1)
	_, err = sr.Seek(sigOffset+4, io.SeekStart)
	if err != nil {
		log.Fatalln(err)
	}
	binary.Read(sr, binary.LittleEndian, &pefile.FileHeader)
	// Print File Header
	fmt.Println("[-----COFF File Header-----]")
	fmt.Printf("[+] Machine Architecture: %#x\n", pefile.FileHeader.Machine)
	fmt.Printf("[+] Number of Sections: %#x\n", pefile.FileHeader.NumberOfSections)
	fmt.Printf("[+] Size of Optional Header: %#x\n", pefile.FileHeader.SizeOfOptionalHeader)
	// Print section names
	fmt.Println("[-----Section Offsets-----]")
	fmt.Printf("[+] Number of Sections Field Offset: %#x\n", sigOffset+6)
	// this is the end of the Signature header (0x7c) + coff (20bytes) + oh32 (224bytes)
	fmt.Printf("[+] Section Table Offset: %#x\n", sigOffset+0xF8)

	var sizeOfOptionalHeader32 = uint16(binary.Size(pe.OptionalHeader32{}))
	var sizeOfOptionalHeader64 = uint16(binary.Size(pe.OptionalHeader64{}))
	var oh32 pe.OptionalHeader32
	var oh64 pe.OptionalHeader64
	switch pefile.FileHeader.SizeOfOptionalHeader {
	case sizeOfOptionalHeader32:
		binary.Read(sr, binary.LittleEndian, &oh32)
	case sizeOfOptionalHeader64:
		binary.Read(sr, binary.LittleEndian, &oh64)
	}
	// Print Optional Header
	fmt.Println("[-----Optional Header-----]")
	fmt.Printf("[+] Entry Point: %#x\n", oh32.AddressOfEntryPoint)
	fmt.Printf("[+] ImageBase: %#x\n", oh32.ImageBase)
	fmt.Printf("[+] Size of Image: %#x\n", oh32.SizeOfImage)
	fmt.Printf("[+] Sections Alignment: %#x\n", oh32.SectionAlignment)
	fmt.Printf("[+] File Alignment: %#x\n", oh32.FileAlignment)
	fmt.Printf("[+] Characteristics: %#x\n", pefile.FileHeader.Characteristics)
	fmt.Printf("[+] Size of Headers: %#x\n", oh32.SizeOfHeaders)
	fmt.Printf("[+] Checksum: %#x\n", oh32.CheckSum)
	fmt.Printf("[+] Machine: %#x\n", pefile.FileHeader.Machine)
	fmt.Printf("[+] Subsystem: %#x\n", oh32.Subsystem)
	fmt.Printf("[+] DLLCharacteristics: %#x\n", oh32.DllCharacteristics)
}
