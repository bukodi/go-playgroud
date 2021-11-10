package jose

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"gopkg.in/square/go-jose.v2"
	_ "gopkg.in/square/go-jose.v2"
	"testing"
)

func TestVault(t *testing.T) {

	pswCallback := func(hint string) string {
		return "Passw0rd"
	}
	// Generate a public/private key pair to use for this example.
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	var data = []byte("Lorem ipsum dolor sit amet")
	encText, err := Encrypt(data, pswCallback, &privateKey.PublicKey)
	if err != nil {
		panic(err)
	}

	data1, err := Decrypt(encText, pswCallback)
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(data, data1) {
		panic("")
	}

	data2, err := DecryptWithBackup(encText, privateKey)
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(data, data2) {
		panic("")
	}
}

type PwsCallback func(hint string) string

func Encrypt(data []byte, pswCallback PwsCallback, backupPublic *ecdsa.PublicKey) (string, error) {

	encrypter, err := jose.NewMultiEncrypter(jose.A128GCM, []jose.Recipient{
		{Algorithm: jose.PBES2_HS256_A128KW, Key: pswCallback("")},
		{Algorithm: jose.ECDH_ES_A256KW, Key: backupPublic},
	}, nil)
	if err != nil {
		return "", err
	}

	object, err := encrypter.Encrypt(data)
	if err != nil {
		return "", err
	}

	return object.FullSerialize(), nil
}

func Decrypt(cipherText string, pswCallback PwsCallback) ([]byte, error) {
	// Parse the serialized, encrypted JWE object. An error would indicate that
	// the given input did not represent a valid message.
	object, err := jose.ParseEncrypted(cipherText)
	if err != nil {
		return nil, err
	}

	_, _, decrypted, err := object.DecryptMulti(pswCallback(""))
	return decrypted, err
}

func DecryptWithBackup(cipherText string, backupKey *ecdsa.PrivateKey) ([]byte, error) {
	// Parse the serialized, encrypted JWE object. An error would indicate that
	// the given input did not represent a valid message.
	object, err := jose.ParseEncrypted(cipherText)
	if err != nil {
		return nil, err
	}

	_, _, decrypted, err := object.DecryptMulti(backupKey)
	return decrypted, err
}

func BenchmarkPasswordDec(b *testing.B) {
	pswCallback := func(hint string) string {
		return "Passw0rd"
	}
	// Generate a public/private key pair to use for this example.
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	var data = []byte("Lorem ipsum dolor sit amet")
	encText, _ := Encrypt(data, pswCallback, &privateKey.PublicKey)
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		Decrypt(encText, pswCallback)
	}
}

func BenchmarkECCDec(b *testing.B) {
	pswCallback := func(hint string) string {
		return "Passw0rd"
	}
	// Generate a public/private key pair to use for this example.
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	var data = []byte("Lorem ipsum dolor sit amet")
	encText, _ := Encrypt(data, pswCallback, &privateKey.PublicKey)

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		DecryptWithBackup(encText, privateKey)
	}
}
