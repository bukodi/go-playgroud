package ecies

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
	"testing"
)

// Based on this description:
// https://medium.com/asecuritysite-when-bob-met-alice/go-public-and-symmetric-key-the-best-of-both-worlds-ecies-180f71eebf59

func TestECIES(t *testing.T) {
	ecKp, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	rnd, err := randomBigInt()
	if err != nil {
		t.Fatal(err)
	}

	Rx, Ry := ecKp.Curve.ScalarMult(ecKp.Curve.Params().Gx, ecKp.Curve.Params().Gy, rnd.Bytes())
	S1x, S1y := ecKp.Curve.ScalarMult(ecKp.X, ecKp.Y, rnd.Bytes())

	S2x, S2y := ecKp.Curve.ScalarMult(Rx, Ry, ecKp.D.Bytes())

	t.Logf("S1 = %+v, %+v", S1x.Bytes(), S1y.Bytes())
	t.Logf("S2 = %+v, %+v", S2x.Bytes(), S2y.Bytes())

}

func randomBigInt() (*big.Int, error) {
	//Max random value, a 130-bits integer, i.e 2^130 - 1
	max := new(big.Int)
	max.Exp(big.NewInt(2), big.NewInt(130), nil).Sub(max, big.NewInt(1))

	//Generate cryptographically strong pseudo-random between 0 - max
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return nil, err
	}
	return n, nil
}
