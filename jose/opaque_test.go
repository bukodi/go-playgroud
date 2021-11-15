package jose

import (
	"bytes"
	"crypto/aes"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"gopkg.in/square/go-jose.v2"
	josecipher "gopkg.in/square/go-jose.v2/cipher"
	"hash"
	"testing"
)

type myCachingPBEDecrypter struct {
	password string
}

func readP2S(p2sIf interface{}) ([]byte, error) {
	if p2sIf == nil {
		return nil, fmt.Errorf("square/go-jose: invalid P2S: must be present")
	}
	p2s := fmt.Sprintf("%s", p2sIf)
	p2sBytes, err := base64.RawURLEncoding.DecodeString(p2s)
	if err != nil {
		return nil, fmt.Errorf("square/go-jose: invalid P2S: invalid base64.: %w", err)
	}
	return p2sBytes, nil
}

func (cpd *myCachingPBEDecrypter) DecryptKey(encryptedKey []byte, header jose.Header) ([]byte, error) {
	alg := jose.KeyAlgorithm(header.Algorithm)
	if !(alg == jose.PBES2_HS256_A128KW || alg == jose.PBES2_HS512_A256KW || alg == jose.PBES2_HS384_A192KW) {
		return nil, jose.ErrUnsupportedAlgorithm
	}

	p2sBytes, err := readP2S(header.ExtraHeaders["p2s"])
	if err != nil {
		return nil, fmt.Errorf("square/go-jose: invalid P2S: invalid base64.: %w", err)
	}
	_ = p2sBytes
	p2cIf := header.ExtraHeaders["p2c"]
	if p2cIf == nil {

	}
	var p2c int
	switch i := p2cIf.(type) {
	case nil:
		return nil, fmt.Errorf("square/go-jose: invalid P2C: must be present")
	case float64:
		p2c = int(i)
	case float32:
		p2c = int(i)
	case int64:
		p2c = int(i)
	default:
		return nil, fmt.Errorf("square/go-jose: invalid P2C: must be present")
	}
	if p2c <= 0 {
		return nil, fmt.Errorf("square/go-jose: invalid P2C: must be a positive integer")
	}

	// salt is UTF8(Alg) || 0x00 || Salt Input
	salt := bytes.Join([][]byte{[]byte(alg), p2sBytes}, []byte{0x00})

	// derive key
	keyLen, h := getPbkdf2Params(alg)
	key := pbkdf2.Key(encryptedKey, salt, p2c, keyLen, h)

	// use AES cipher with derived key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	cek, err := josecipher.KeyUnwrap(block, encryptedKey)
	if err != nil {
		return nil, err
	}
	return cek, nil

}

// getPbkdf2Params returns the key length and hash function used in
// pbkdf2.Key.
func getPbkdf2Params(alg jose.KeyAlgorithm) (int, func() hash.Hash) {
	switch alg {
	case jose.PBES2_HS256_A128KW:
		return 16, sha256.New
	case jose.PBES2_HS384_A192KW:
		return 24, sha512.New384
	case jose.PBES2_HS512_A256KW:
		return 32, sha512.New
	default:
		panic("invalid algorithm")
	}
}

func TestJWEPBEForOpaque(t *testing.T) {

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

	decrypted1, err := object.Decrypt("Passw0rd")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Decrypted: %s\n", string(decrypted1))

	mcpd := &myCachingPBEDecrypter{"Passw0rd"}
	decrypted2, err := object.Decrypt(mcpd)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Decrypted: %s\n", string(decrypted2))
}
