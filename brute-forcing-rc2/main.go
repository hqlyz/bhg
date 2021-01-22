package main

import (
	"brute-forcing-rc2/rc2"
	"crypto/cipher"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"regexp"
	"sync"

	luhn "github.com/joeljunstrom/go-luhn"
)

var numeric = regexp.MustCompile(`^\d{8}$`)

type CryptoData struct {
	block cipher.Block
	key   []byte
}

func generate(start, stop uint64, out chan<- *CryptoData, done <-chan struct{}, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := start; i <= stop; i++ {
			key := make([]byte, 8)
			select {
			case <-done:
				return
			default:
				binary.BigEndian.PutUint64(key, i)
				block, err := rc2.New(key[3:], 40)
				if err != nil {
					log.Fatalf("create rc2 cipher failed: %v\n", err)
				}
				data := &CryptoData{
					block: block,
					key:   key[3:],
				}
				out <- data
			}
		}
	}()
}

func decrypt(ciphertext []byte, in <-chan *CryptoData, done chan struct{}, wg *sync.WaitGroup) {
	size := rc2.BlockSize
	plainText := make([]byte, len(ciphertext))
	wg.Add(1)
	go func() {
		defer wg.Done()
		for data := range in {
			select {
			case <-done:
				return
			default:
				data.block.Decrypt(plainText[:size], ciphertext[:size])
				if numeric.Match(plainText[:size]) {
					data.block.Decrypt(plainText[size:], ciphertext[size:])
					if luhn.Valid(string(plainText)) && numeric.Match(plainText[size:]) {
						fmt.Printf("found %s with key [%x]\n", string(plainText), data.key)
						close(done)
						return
					}
				}
			}
		}
	}()
}

func main() {
	ciphertext, err := hex.DecodeString("0986f2cc1ebdc5c2e25d04a136fa1a6b")
	if err != nil {
		log.Fatalln(err)
	}

	var (
		min, max, prods    = uint64(0x00000000), uint64(0xffffffff), uint64(75)
		generateWg, workWg sync.WaitGroup
		step               = (max - min) / prods
		start, end         = min, min + step
		workCount          = 100
		done               = make(chan struct{})
		work               = make(chan *CryptoData, workCount)
	)

	for i := uint64(0); i <= prods; i++ {
		if end > max {
			end = max
		}
		generate(start, end, work, done, &generateWg)

		start += step
		end += step
	}

	log.Println("Producers started!")
	log.Println("Starting consumers...")
	for i := 0; i < workCount; i++ {
		decrypt(ciphertext, work, done, &workWg)
	}
	log.Println("Consumers started!")
	log.Println("Now we wait...")
	generateWg.Wait()
	close(work)
	workWg.Wait()
	log.Println("Brute-force complete")
}
