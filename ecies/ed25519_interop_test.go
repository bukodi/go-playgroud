package ecies

import (
	stdEd25519 "crypto/ed25519"
	"crypto/rand"
	xEd25519 "golang.org/x/crypto/ed25519"
	"testing"
)

func TestEd25519Interoperability(t *testing.T) {
	msg := []byte("Hello")

	stdEdPub, stdEdPriv, _ := stdEd25519.GenerateKey(rand.Reader)
	xEdPub, xEdPriv, _ := xEd25519.GenerateKey(rand.Reader)

	stdImpPriv := stdEd25519.NewKeyFromSeed(xEdPriv.Seed())
	xImpPriv := xEd25519.NewKeyFromSeed(stdEdPriv.Seed())
	xImpPub, ok := xImpPriv.Public().(xEd25519.PublicKey)
	if !ok {
		t.FailNow()
	}

	_ = stdEdPub
	_ = stdEdPriv
	_ = xEdPub
	_ = xEdPriv
	_ = stdImpPriv
	_ = xImpPriv

	if !stdEd25519.Verify(stdEdPub, msg, stdEd25519.Sign(stdEdPriv, msg)) {
		t.FailNow()
	}
	if !xEd25519.Verify(xImpPub, msg, xEd25519.Sign(stdEdPriv, msg)) {
		t.FailNow()
	}
}
