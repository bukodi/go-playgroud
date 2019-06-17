package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"testing"
)

func TestCompositeEC(t *testing.T) {
	privKey1, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	privKey2, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	privKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	privKey.D.Add(privKey1.D, privKey2.D)
	privKey.PublicKey.X.Add(privKey1.PublicKey.X, privKey2.PublicKey.X)
	privKey.PublicKey.Y.Add(privKey1.PublicKey.Y, privKey2.PublicKey.Y)

	msg := "hello, world"
	hash := sha256.Sum256([]byte(msg))

	r, s, err := ecdsa.Sign(rand.Reader, privKey, hash[:])
	if err != nil {
		panic(err)
	}
	fmt.Printf("signature: (0x%x, 0x%x)\n", r, s)

	valid := ecdsa.Verify(&privKey.PublicKey, hash[:], r, s)
	fmt.Println("signature verified:", valid)
}

func TestSimpleECSign(t *testing.T) {
	privKey1, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	msg := "hello, world"
	hash := sha256.Sum256([]byte(msg))

	r, s, err := ecdsa.Sign(rand.Reader, privKey1, hash[:])
	if err != nil {
		panic(err)
	}
	fmt.Printf("signature: (0x%x, 0x%x)\n", r, s)

	valid := ecdsa.Verify(&privKey1.PublicKey, hash[:], r, s)
	fmt.Println("signature verified:", valid)
}
