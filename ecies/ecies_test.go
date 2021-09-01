package ecies

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"io"
	"math/big"
	"testing"
)

// Based on this description:
// https://medium.com/asecuritysite-when-bob-met-alice/go-public-and-symmetric-key-the-best-of-both-worlds-ecies-180f71eebf59

func TestECIES_P256(t *testing.T) {
	ecKp, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	rnd, err := randFieldElement(ecKp.Curve, rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	Rx, Ry := ecKp.Curve.ScalarMult(ecKp.Curve.Params().Gx, ecKp.Curve.Params().Gy, rnd.Bytes())
	S1x, S1y := ecKp.Curve.ScalarMult(ecKp.X, ecKp.Y, rnd.Bytes())

	S2x, S2y := ecKp.Curve.ScalarMult(Rx, Ry, ecKp.D.Bytes())

	t.Logf("S1 = %+v, %+v", S1x.Bytes(), S1y.Bytes())
	t.Logf("S2 = %+v, %+v", S2x.Bytes(), S2y.Bytes())
	if S1x.Cmp(S2x) != 0 || S1y.Cmp(S2y) != 0 {
		t.Fatal("shared secrets aren't equals")
	}

}

func TestECIES_ed25519(t *testing.T) {
	edPub, edPriv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	signature := ed25519.Sign(edPriv, []byte("Hello"))
	t.Logf("Signature: %v", signature)

	ok := ed25519.Verify(edPub, []byte("Hello"), signature)
	t.Logf("Signature valid: %v", ok)

}

var one = new(big.Int).SetInt64(1)

// randFieldElement returns a random element of the field underlying the given
// curve using the procedure given in [NSA] A.2.1.
func randFieldElement(c elliptic.Curve, rand io.Reader) (k *big.Int, err error) {
	params := c.Params()
	b := make([]byte, params.BitSize/8+8)
	_, err = io.ReadFull(rand, b)
	if err != nil {
		return
	}

	k = new(big.Int).SetBytes(b)
	n := new(big.Int).Sub(params.N, one)
	k.Mod(k, n)
	k.Add(k, one)
	return
}
