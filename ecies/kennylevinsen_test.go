package ecies

import (
	"bytes"
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"testing"

	_ "github.com/kennylevinsen/ecies"
	"golang.org/x/crypto/curve25519"
)

// original : https://github.com/kennylevinsen/ecies

func Encrypt(plainText, publicKey []byte) ([]byte, error) {
	var r, R, S, K_B [32]byte

	if _, err := rand.Read(r[:]); err != nil {
		return nil, err
	}
	r[0] &= 248
	r[31] &= 127
	r[31] |= 64

	copy(K_B[:], publicKey)

	curve25519.ScalarBaseMult(&R, &r)
	curve25519.ScalarMult(&S, &r, &K_B)
	k_E := sha512.Sum512(S[:])

	cipherText := make([]byte, 32+len(plainText))
	copy(cipherText[:32], R[:])
	for i := 0; i < len(plainText); i++ {
		cipherText[32+i] = plainText[i] ^ k_E[i]
	}

	return cipherText, nil
}

func Decrypt(cipherText, privateKey []byte) ([]byte, error) {
	var R, S, k_B [32]byte
	copy(R[:], cipherText[:32])
	copy(k_B[:], privateKey)

	curve25519.ScalarMult(&S, &k_B, &R)

	k_E := sha512.Sum512(S[:])

	plainText := make([]byte, len(cipherText)-32)
	for i := 0; i < len(plainText); i++ {
		plainText[i] = cipherText[32+i] ^ k_E[i]
	}

	return plainText, nil
}

func TestEncrypt(t *testing.T) {
	var privKey, pubKey [32]byte

	if _, err := rand.Read(privKey[:]); err != nil {
		t.Fatalf("could not generate privkey: %v", err)
	}

	privKey[0] &= 248
	privKey[31] &= 127
	privKey[31] |= 64

	curve25519.ScalarBaseMult(&pubKey, &privKey)

	plainText := []byte("Hej, mit navn er Per. Jeg kan godt lide ost.")

	cipherText, err := Encrypt(plainText, pubKey[:])
	if err != nil {
		t.Fatalf("got error: %v", err)
	}

	plainText2, err := Decrypt(cipherText, privKey[:])
	if err != nil {
		t.Fatalf("got error: %v", err)
	}

	if bytes.Compare(plainText, plainText2) != 0 {
		t.Fatalf("result did not match:\nGot:\n%s\nExpected:\n%s\n", hex.Dump(plainText2), hex.Dump(plainText))
	}
}

func TestDecrypt(t *testing.T) {
	privKey := []byte{0xc8, 0x06, 0x43, 0x9d, 0xc9, 0xd2, 0xc4, 0x76, 0xff, 0xed, 0x8f, 0x25, 0x80, 0xc0, 0x88, 0x8d, 0x58, 0xab, 0x40, 0x6b, 0xf7, 0xae, 0x36, 0x98, 0x87, 0x90, 0x21, 0xb9, 0x6b, 0xb4, 0xbf, 0x59}
	cipherText := []byte{0xDA, 0xBF, 0x5E, 0x74, 0xB8, 0x43, 0x09, 0xBC, 0x5B, 0x9E, 0xC9, 0x69, 0x79, 0x02, 0x39, 0xA8, 0x71, 0xD5, 0xC6, 0xC5, 0xE9, 0x9C, 0xC3, 0x04, 0xE5, 0x87, 0x58, 0xBC, 0xD8, 0x5F, 0x8F, 0x50, 0x1D, 0x67, 0xF4, 0x10, 0xDA, 0x39, 0xD2, 0xFC, 0x3F, 0x87, 0x85, 0xE4, 0x84, 0xE1, 0x61, 0xFB, 0xA0, 0x45, 0x0A, 0x60, 0x49, 0x2A, 0x4F, 0x91, 0x97, 0x9D, 0xC7, 0xFF}
	plainText := []byte("Hello there, my name is Paul")
	plainText2, err := Decrypt(cipherText, privKey)
	if err != nil {
		t.Fatalf("got error: %v", err)
	}

	if bytes.Compare(plainText, plainText2) != 0 {
		t.Fatalf("result did not match:\nGot:\n%s\nExpected:\n%s\n", hex.Dump(plainText2), hex.Dump(plainText))
	}
}
