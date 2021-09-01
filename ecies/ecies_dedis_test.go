package ecies

import (
	"bytes"
	"go.dedis.ch/kyber/v3/encrypt/ecies"
	"go.dedis.ch/kyber/v3/group/edwards25519"
	"go.dedis.ch/kyber/v3/util/random"
	"testing"
)

func TestECIES_DEDIS(t *testing.T) {

	message := []byte("Hello")

	suite := edwards25519.NewBlakeSHA256Ed25519()

	private := suite.Scalar().Pick(random.New())

	public := suite.Point().Mul(private, nil)

	ciphertext, _ := ecies.Encrypt(suite, public, message, suite.Hash)

	t.Logf("Text input: %s\n\n", message)
	t.Logf("Private key: %s\n", private)
	t.Logf("Public key: %s\n\n", public)

	t.Logf("Cipher: %x", ciphertext)

	plaintext, _ := ecies.Decrypt(suite, private, ciphertext, suite.Hash)
	t.Logf("Decipher: %s", plaintext)
	if !bytes.Equal(plaintext, message) {
		t.Fatalf("Wrog decryption: %s", plaintext)
	}
}
