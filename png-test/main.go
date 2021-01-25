package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
)

type Header struct {
	Header uint64
}

type Chunk struct {
	Size uint32
	Type uint32
	Data []byte
	CRC  uint32
}

type MetaChunk struct {
	chk Chunk
}

//PreProcessImage reads to buffer from file handle
func PreProcessImage(dat *os.File) (*bytes.Reader, error) {
	stats, err := dat.Stat()
	if err != nil {
		return nil, err
	}
	var size = stats.Size()
	b := make([]byte, size)
	bufR := bufio.NewReader(dat)
	_, err = bufR.Read(b)
	bReader := bytes.NewReader(b)
	bReader.Seek(0, io.SeekCurrent)
	return bReader, err
}

func (mc *MetaChunk) validate(b *bytes.Reader) {
	var header Header
	if err := binary.Read(b, binary.BigEndian, &header.Header); err != nil {
		log.Fatal(err)
	}
	bArr := make([]byte, 8)
	binary.BigEndian.PutUint64(bArr, header.Header)
	if string(bArr[1:4]) != "PNG" {
		log.Fatal("Provided file is not a valid PNG format")
	} else {
		fmt.Println("Valid PNG so let us continue!")
	}
}

func main() {

}
