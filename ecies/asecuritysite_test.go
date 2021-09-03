package ecies

import (
	"bytes"
	"golang.org/x/crypto/curve25519"
	"math/rand"
	"testing"
	"time"
)

func TestECDH(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	var privateKey [32]byte
	for i := range privateKey[:] {
		privateKey[i] = byte(rand.Intn(256))
	}

	var publicKey [32]byte
	curve25519.ScalarBaseMult(&publicKey, &privateKey)

	t.Logf("\nAlice Private key (a):\t%x\n", privateKey)
	t.Logf("\nAlice Public key point (x co-ord):\t%x\n", publicKey)

	var privateKey2 [32]byte
	for i := range privateKey[:] {
		privateKey2[i] = byte(rand.Intn(256))
	}

	var publicKey2 [32]byte
	curve25519.ScalarBaseMult(&publicKey2, &privateKey2)

	var out1, out2 [32]byte

	t.Logf("\nBob Private key (b):\t%x\n", privateKey2)
	t.Logf("\nBob Public key point (x co-ord):\t%x\n", publicKey2)

	curve25519.ScalarMult(&out1, &privateKey, &publicKey2)
	curve25519.ScalarMult(&out2, &privateKey2, &publicKey)

	t.Logf("\nShared key (Alice):\t%x\n", out1)
	t.Logf("\nShared key (Bob):\t%x\n", out2)

	if !bytes.Equal(out1[:], out2[:]) {
		t.FailNow()
	}

}
