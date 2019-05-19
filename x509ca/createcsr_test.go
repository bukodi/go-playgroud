package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"testing"
)

type MySignerKey struct {
	privKey *rsa.PrivateKey
}

func (mySigner *MySignerKey) Public() crypto.PublicKey {
	return mySigner.privKey.Public()
}

func (mySigner *MySignerKey) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	return mySigner.privKey.Sign(rand, digest, opts)
}

func TestCreateCsr(t *testing.T) {
	// step: generate a keypair
	keys, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Error(err)
	}

	mySigner := MySignerKey{
		privKey: keys,
	}

	names := pkix.Name{
		CommonName:   "Cica Mica",
		Organization: []string{"Test"},
		Country:      []string{"HU"},
	}

	csrPem, err := createCsr(names, &mySigner)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%s", csrPem)

}

//
// createCertificateAuthority generates a certificate authority request ready to be signed
//
func createCsr(names pkix.Name, keys crypto.Signer) (string, error) {

	// step: generate a csr template
	var csrTemplate = x509.CertificateRequest{
		Subject:            names,
		SignatureAlgorithm: x509.SHA512WithRSA,
	}
	// step: generate the csr request
	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &csrTemplate, keys)
	if err != nil {
		return "", err
	}
	csrPem := pem.EncodeToMemory(&pem.Block{
		Type: "CERTIFICATE REQUEST", Bytes: csrBytes,
	})
	return string(csrPem), nil
}
