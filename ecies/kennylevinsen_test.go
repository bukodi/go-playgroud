package ecies

import (
	"bytes"
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"golang.org/x/crypto/ed25519"
	"testing"

	_ "github.com/kennylevinsen/ecies"
	"golang.org/x/crypto/curve25519"
)

// original : https://github.com/kennylevinsen/ecies

func Encrypt(plainText, publicKey []byte) ([]byte, error) {
	var rndScalar [32]byte

	if _, err := rand.Read(rndScalar[:]); err != nil {
		return nil, err
	}
	rndScalar[0] &= 248
	rndScalar[31] &= 127
	rndScalar[31] |= 64

	rndPoint, err := curve25519.X25519(rndScalar[:], curve25519.Basepoint)
	if err != nil {
		return nil, err
	}
	sharedSecret, err := curve25519.X25519(rndScalar[:], publicKey)
	if err != nil {
		return nil, err
	}
	digest := sha512.Sum512(sharedSecret)

	cipherText := make([]byte, 32+len(plainText))
	copy(cipherText[:32], rndPoint[:])
	for i := 0; i < len(plainText); i++ {
		cipherText[32+i] = plainText[i] ^ digest[i]
	}

	return cipherText, nil
}

func Decrypt(cipherText, privateKey []byte) ([]byte, error) {
	sharedSecret, err := curve25519.X25519(privateKey, cipherText[:32])
	if err != nil {
		return nil, err
	}

	digest := sha512.Sum512(sharedSecret[:])

	plainText := make([]byte, len(cipherText)-32)
	for i := 0; i < len(plainText); i++ {
		plainText[i] = cipherText[32+i] ^ digest[i]
	}

	return plainText, nil
}

func TestBoth(t *testing.T) {

	t.Run("keygen", TestKeyGen)
	t.Run("normal", TestEncrypt)
	t.Run("failed", TestEncrypt2)
}

func TestKeyGen(t *testing.T) {
	edPub, edPriv, _ := ed25519.GenerateKey(fakeRand)
	privKey := scalarFromSeed([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	pubKey, _ := curve25519.X25519(privKey, curve25519.Basepoint)
	var pubKey2 [32]byte
	var privTmp [32]byte
	copy(privTmp[:], privKey)
	curve25519.ScalarBaseMult(&pubKey2, &privTmp)
	t.Logf("privKey=%+v", privKey)
	t.Logf("pubKey=%+v", pubKey)
	t.Logf("pubKey2=%+v", pubKey2)

	t.Logf("edPriv=%+v", scalarFromSeed(edPriv.Seed()))
	t.Logf("edPub=%+v", edPub)

}

func TestEncrypt(t *testing.T) {
	privKey := []byte{175, 34, 224, 240, 87, 185, 220, 205, 75, 27, 229, 206, 119, 226, 231, 213, 87, 181, 121, 112, 181, 38, 122, 144, 245, 121, 96, 146, 74, 135, 241, 6}
	//privKey := []byte{57,132,30,233,176,214,200,247,124,41,25,23,100,18,99,193,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,15}
	//privKey := []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}
	//privKey[0] &= 248
	//privKey[31] &= 127
	//privKey[31] |= 64

	t.Logf("privKey=%+v", privKey)
	pubKey, _ := curve25519.X25519(privKey, curve25519.Basepoint)
	t.Logf("pubKey=%+v", pubKey)

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

func TestEncrypt2(t *testing.T) {
	edPub, edPriv, _ := ed25519.GenerateKey(fakeRand)
	h := sha512.Sum512(edPriv.Seed())
	privKey := setBytesWithClamping(h[:32])
	t.Logf("privKey=%+v", privKey)
	t.Logf("pubKey=%+v", edPub)

	plainText := []byte("Hej, mit navn er Per. Jeg kan godt lide ost.")

	cipherText, err := Encrypt(plainText, edPub)
	if err != nil {
		t.Fatalf("got error: %v", err)
	}

	plainText2, err := Decrypt(cipherText, privKey)
	if err != nil {
		t.Fatalf("got error: %v", err)
	}

	if bytes.Compare(plainText, plainText2) != 0 {
		t.Fatalf("result did not match:\nGot:\n%s\nExpected:\n%s\n", hex.Dump(plainText2), hex.Dump(plainText))
	}
}

func scalarFromSeed(edPriv ed25519.PrivateKey) []byte {
	h := sha512.Sum512(edPriv.Seed())
	return setBytesWithClamping(h[:32])
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
