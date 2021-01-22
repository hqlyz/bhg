package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"log"
)

func main() {
	var (
		err                                              error
		publicKey                                        *rsa.PublicKey
		privateKey                                       *rsa.PrivateKey
		plainText, cipherText, label, signature, message []byte
	)

	if privateKey, err = rsa.GenerateKey(rand.Reader, 2048); err != nil {
		log.Fatalf("rsa generate key failed: %v\n", err)
	}
	publicKey = &privateKey.PublicKey
	label = []byte("")
	message = []byte("Some super secret message, maybe a session key even")

	if cipherText, err = rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, message, label); err != nil {
		log.Fatalf("rsa encrypt failed: %v\n", err)
	}
	fmt.Printf("ciphertext: %x\n", string(cipherText))

	if plainText, err = rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, cipherText, label); err != nil {
		log.Fatalf("rsa decrypt failed: %v\n", err)
	}
	fmt.Printf("plaintext: %s\n", plainText)

	// Signature validate
	h := sha256.New()
	h.Write(message)
	if signature, err = rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, h.Sum(nil), nil); err != nil {
		log.Fatalf("rsa signPSS failed: %v\n", err)
	}
	fmt.Printf("signature: %x\n", signature)

	if err = rsa.VerifyPSS(publicKey, crypto.SHA256, h.Sum(nil), signature, nil); err != nil {
		log.Fatalf("rsa verifyPSS failed: %v\n", err)
	}
	fmt.Println("signature verified")
}
