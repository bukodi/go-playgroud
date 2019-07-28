package certutil

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"github.com/fullsailor/pkcs7"
	"math/big"
	"testing"
	"time"
)

func TestWrappedImport(t *testing.T) {
	wCert, wKey, err := CreateWrappingKey()
	if err != nil {
		t.Fatal(err)
	}

	encMsg, err := ExportProtectedData([]byte("Minden cica aranyos"), wCert)
	if err != nil {
		t.Fatal(err)
	}

	data, err := ImportProtectedData(encMsg, wCert, wKey)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Decrypted data: %s\n", data)
}

func CreateWrappingKey() (wrappingCertPEM string, wrappingPrivKey *rsa.PrivateKey, err error) {
	wrappingKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return "", nil, err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "WrappingKey",
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24 * 180),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	wrappingCerDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &wrappingKey.PublicKey, wrappingKey)
	if err != nil {
		return "", nil, err
	}
	b := pem.Block{
		Type:  "CERTIFICATE",
		Bytes: wrappingCerDER,
	}
	certPEM := pem.EncodeToMemory(&b)
	return string(certPEM), wrappingKey, nil
}

func pemToCert(certPEM string) (*x509.Certificate, error) {
	parsedPEM, _ := pem.Decode([]byte(certPEM))
	if parsedPEM.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("Invalid PEM type: %s (Expected type: %s)", parsedPEM.Type, "CERTIFICATE")
	}

	cert, err := x509.ParseCertificate(parsedPEM.Bytes)
	if err != nil {
		return nil, err
	}
	return cert, nil
}

func ExportProtectedData(data []byte, wrappingCertPEM string) (p7PEM string, err error) {
	wrappingCert, err := pemToCert(wrappingCertPEM)
	if err != nil {
		return "", err
	}

	encP7, err := pkcs7.Encrypt(data, []*x509.Certificate{wrappingCert})
	if err != nil {
		return "", err
	}
	b := pem.Block{
		Type:  "PKCS7",
		Bytes: encP7,
	}
	pemBlock := pem.EncodeToMemory(&b)
	return string(pemBlock), nil
}

func ImportProtectedData(p7PEM string, wrappingCertPEM string, wrappingPrivKey *rsa.PrivateKey) (data []byte, err error) {
	wrappingCert, err := pemToCert(wrappingCertPEM)
	if err != nil {
		return nil, err
	}

	p7ParsedPEM, _ := pem.Decode([]byte(p7PEM))
	if p7ParsedPEM.Type != "PKCS7" {
		return nil, fmt.Errorf("Invalid PEM type: %s (Expected type: %s)", p7ParsedPEM.Type, "PKCS7")
	}

	p7MsgObj, err := pkcs7.Parse(p7ParsedPEM.Bytes)
	if err != nil {
		return nil, err
	}

	data, err = p7MsgObj.Decrypt(wrappingCert, wrappingPrivKey)
	if err != nil {
		return nil, err
	}
	return data, nil
}
