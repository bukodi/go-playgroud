package jose

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"sync"
	"testing"
)
import "github.com/go-jose/go-jose/v3"

func TestVault(t *testing.T) {

	pswCallback := func(hint string) []byte {
		return []byte("Passw0rd")
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

func Encrypt(data []byte, pswCallback PswCallback, backupPublic *ecdsa.PublicKey) (string, error) {

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

func Decrypt(cipherText string, pswCallback PswCallback) ([]byte, error) {
	// Parse the serialized, encrypted JWE object. An error would indicate that
	// the given input did not represent a valid message.
	object, err := jose.ParseEncrypted(cipherText)
	if err != nil {
		return nil, err
	}

	_, _, decrypted, err := object.DecryptMulti(pswCallback(""))
	return decrypted, err
}

func DecryptWithMemoizerPBE(cipherText string, pswCallback PswCallback) ([]byte, error) {
	// Parse the serialized, encrypted JWE object. An error would indicate that
	// the given input did not represent a valid message.
	object, err := jose.ParseEncrypted(cipherText)
	if err != nil {
		return nil, err
	}

	mcpd := &myCachingPBEDecrypter{pswCallback}

	_, _, decrypted, err := object.DecryptMulti(mcpd)
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
	pswCallback := func(hint string) []byte {
		return []byte("Passw0rd")
	}
	// Generate a public/private key pair to use for this example.
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	var data = []byte("Lorem ipsum dolor sit amet")
	encText, _ := Encrypt(data, pswCallback, &privateKey.PublicKey)
	// run the Fib function b.N times
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Decrypt(encText, pswCallback)
	}
}

func BenchmarkCacheablePasswordDec(b *testing.B) {
	pswCallback := func(hint string) []byte {
		return []byte("Passw0rd")
	}
	// Generate a public/private key pair to use for this example.
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	var data = []byte("Lorem ipsum dolor sit amet")
	encText, _ := Encrypt(data, pswCallback, &privateKey.PublicKey)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := DecryptWithMemoizerPBE(encText, pswCallback)
		if err != nil {
			b.Fatalf("%+v", err)
		}
	}
}

func TestMemoizerPBE(t *testing.T) {
	pswCallback := func(hint string) []byte {
		return []byte("Passw0rd")
	}
	// Generate a public/private key pair to use for this example.
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	var data = []byte("Lorem ipsum dolor sit amet")
	encText, _ := Encrypt(data, pswCallback, &privateKey.PublicKey)

	var wg sync.WaitGroup
	wg.Add(10)
	for m := 0; m < 10; m++ {
		go func() {
			for n := 0; n < 10; n++ {
				plaintext, err := DecryptWithMemoizerPBE(encText, pswCallback)
				if err != nil {
					t.Fatalf("%+v", err)
				} else {
					t.Logf("%s", plaintext)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkMemoizerPBE(t *testing.B) {
	pswCallback := func(hint string) []byte {
		return []byte("Passw0rd")
	}
	// Generate a public/private key pair to use for this example.
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	var data = []byte("Lorem ipsum dolor sit amet")
	encText, _ := Encrypt(data, pswCallback, &privateKey.PublicKey)

	for n := 0; n < t.N; n++ {
		plaintext, err := DecryptWithMemoizerPBE(encText, pswCallback)
		if err != nil {
			t.Fatalf("%+v", err)
		} else {
			t.Logf("%s", plaintext)
		}
	}
}

func BenchmarkECCDec(b *testing.B) {
	pswCallback := func(hint string) []byte {
		return []byte("Passw0rd")
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
