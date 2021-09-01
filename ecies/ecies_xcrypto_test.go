package ecies

import (
	"bytes"
	"crypto/rand"
	"crypto/sha512"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/ed25519"
	"io"
	"testing"
)

func TestECIES_xcrypto(t *testing.T) {
	edPub, edPriv, err := ed25519.GenerateKey(rand.Reader)

	rnd := make([]byte, 32)
	_, err = io.ReadFull(rand.Reader, rnd)
	if err != nil {
		t.Fatal(err)
	}

	R, err := curve25519.X25519(rnd, curve25519.Basepoint)
	if err != nil {
		t.Fatal(err)
	}
	S1, err := curve25519.X25519(rnd, edPub)
	if err != nil {
		t.Fatal(err)
	}

	digest := sha512.Sum512(edPriv.Seed())
	digest[0] &= 248
	digest[31] &= 127
	digest[31] |= 64

	var hBytes [32]byte
	copy(hBytes[:], digest[:])

	S2, err := curve25519.X25519(hBytes[:], R)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("S1 = %+v", S1)
	t.Logf("S1 = %+v", S2)
	if !bytes.Equal(S1, S2) {
		t.Fatal("shared secrets aren't equals")
	}

}
