package jose

import (
	"fmt"
	"gopkg.in/square/go-jose.v2"
	"testing"
)

type myCachingPBEDecrypter struct {
	password string
}

func (cpd *myCachingPBEDecrypter) DecryptKey(encryptedKey []byte, header jose.Header) ([]byte, error) {
	return nil, nil
}

func TestJWEPBEForOpaqe(t *testing.T) {

	encrypter, err := jose.NewEncrypter(jose.A128GCM, jose.Recipient{Algorithm: jose.PBES2_HS256_A128KW, Key: "Passw0rd"}, nil)
	if err != nil {
		panic(err)
	}

	var plaintext = []byte("Lorem ipsum dolor sit amet")
	object, err := encrypter.Encrypt(plaintext)
	if err != nil {
		panic(err)
	}

	serialized := object.FullSerialize()
	fmt.Printf("Serialized message: %s\n", serialized)

	object, err = jose.ParseEncrypted(serialized)
	if err != nil {
		panic(err)
	}

	object.GetAuthData()

	decrypted, err := object.Decrypt("Passw0rd")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Decrypted: %s\n", string(decrypted))
}
