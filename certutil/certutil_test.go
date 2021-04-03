package certutil

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"github.com/fullsailor/pkcs7"
	"math/big"
	"strings"
	"testing"
	"time"
)

func TestWrappedExportImport(t *testing.T) {
	wCert, wKey, err := CreateWrappingKey()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Wrapping cert: \n%s\n", wCert)

	encMsg, err := ExportProtectedData([]byte("Minden cica aranyos"), wCert)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Encrypted msg: \n%s\n", encMsg)

	data, err := ImportProtectedData(encMsg, wCert, wKey)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Decrypted data: %s\n", data)
}

const wrappingPrivkey = `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQCtT0yHKhhJxIXezSn0VpH6YQ0dR4yjJV0hh3WHyDVMft/BAwpo
dMA5YMVSwOyfP/dLbphSZBQ3hx9kvh/mQodi88EjfrkZuB9GRHm8tvFSMxZpzihX
we89/0k6rrhnL/JMIi5EKjalX7/g8A54/iUc3qANwsHkjxpqs5hF+/xS7QIDAQAB
AoGAJYZb/ZAhvictDTbC9MKhzOKiokpRqyl15sKbshCpaay26eWspT1SOZo4X1ii
Y0DIXK7PkFdXAUbihz+bTJHvGiR7WTUb6d1R2pdEcQ9PkRw8UCdpVmy0yHdXn3RH
jaOspuHjGArjSw4yXmE+YBoRMfmwK8eQHhKJPkCeuVolnAECQQDcKDSuCqs7PZbI
suIid5rRLgpRpoNX7Hb7iUu4GMXa2iMcmjqilHAnkjAsxnA9xRbRT8CepAqObz7V
BJjb7C9tAkEAyYaSp3Su6FaOTdy2paOZgqTt3rHbhXu9BIq8DcUVo0c8o/pCEfun
0wQBctdvUuputY9q/EjCZLMYWgfuok0BgQJAacYLYV8Q27qAelavnP237S5gmNAW
pKSgmcNiFMYaMUbCvdg/uaL1q51p9ek1Pgg2KURW+CE6TSqXloLJ6ke0lQJBAIDo
HO0SnUMvAq3ZjdecM3i3CRUXDcVvpjO+jpX8SF/+FFMLpLGeGPbIrCiKl8IhBRoM
GDIyVG5XhS8pvcKBoYECQQDOthuYaJlkgdlO79DG+c0kAKCcR3ZOvYbmaih86o2J
tOi7brbCI317tLpvqW5/QzfvraIqVRjgftxX7CaxCVtP
-----END RSA PRIVATE KEY-----`

const wrappingCert = `-----BEGIN CERTIFICATE-----
MIIBwjCCASugAwIBAgIBATANBgkqhkiG9w0BAQsFADAWMRQwEgYDVQQDEwtXcmFw
cGluZ0tleTAeFw0xOTA3MjgyMTA3NTdaFw0yMDAxMjQyMTA3NTdaMBYxFDASBgNV
BAMTC1dyYXBwaW5nS2V5MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCtT0yH
KhhJxIXezSn0VpH6YQ0dR4yjJV0hh3WHyDVMft/BAwpodMA5YMVSwOyfP/dLbphS
ZBQ3hx9kvh/mQodi88EjfrkZuB9GRHm8tvFSMxZpzihXwe89/0k6rrhnL/JMIi5E
KjalX7/g8A54/iUc3qANwsHkjxpqs5hF+/xS7QIDAQABoyAwHjAOBgNVHQ8BAf8E
BAMCAqQwDAYDVR0TAQH/BAIwADANBgkqhkiG9w0BAQsFAAOBgQBHG0lacpYH1jgl
6dcrr1zA+JtZ2iSst+2msmIArgmqKWFr8PihyO6tIVOk90IduZOsvUgdbmNRRXFw
J/ifOY1nLqb08Je0xSPiV87KnruWwNZXxygi6yMYBox4OChZU/9H967+kNdNuPDj
9IH5xNxFAsCihD/gmaBD8iqdSQf0Fg==
-----END CERTIFICATE-----`

const encryptedMessage = `-----BEGIN PKCS7-----
MIAGCSqGSIb3DQEHA6CAMIACAQAxgbUwgbICAQAwGzAWMRQwEgYDVQQDEwtXcmFw
cGluZ0tleQIBATANBgkqhkiG9w0BAQEFAASBgIr0x4tAJD+h8DUPn1sT4edVNI6n
v2rIuVpbtBBCm3rB3ZUyNxyv2KvK4TZd1Ic0NDZ1q5DWHHIvsuZjCLw5l4k5NM9e
sDB5sCKCYeRJhbIJnfD2vykMUYCXMzREvxewRb0PglB4O7Bh+HOsDcjRcO3p+FFm
n9cUmTpxYQ31kbgbMIAGCSqGSIb3DQEHATAdBglghkgBZQMEAQIEEI504eGznGwa
+fFUv+uEzfmggAQg303lCxENMjeTt92WzW9X63HslmePMOCDOSk8C4G+V80AAAAA
AAAAAAAA
-----END PKCS7-----`

func TestImportWithKnownKey(t *testing.T) {
	parsedPEM, _ := pem.Decode([]byte(wrappingPrivkey))
	wKey, _ := x509.ParsePKCS1PrivateKey(parsedPEM.Bytes)

	data, err := ImportProtectedData(encryptedMessage, wrappingCert, wKey)
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

	pkcs7.ContentEncryptionAlgorithm = pkcs7.EncryptionAlgorithmAES128GCM
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

	pkcs7.ContentEncryptionAlgorithm = pkcs7.EncryptionAlgorithmAES128GCM
	data, err = p7MsgObj.Decrypt(wrappingCert, wrappingPrivKey)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func TestCalcHash(t *testing.T) {
	sum := sha256.Sum256([]byte("Minden cica aranyos."))
	fmt.Printf("% x", sum)
}

func TestParseM5StickKey(t *testing.T) {
	m5KeyPEM := `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIC30UVjPjLFARva1SykxA0cEW3AwtF39IHh/ixrYWSlQoAoGCCqGSM49
AwEHoUQDQgAEoopFL6hGLB+Z9NdmU32JtpTE2lawJpIBGK6ouNXURAujeC8ctyQr
DnN1VtqjyszqkxeDSdR8Bi5NXN7eGsX1dw==
-----END EC PRIVATE KEY-----`

	m5KeyParsedPEM, _ := pem.Decode([]byte(m5KeyPEM))
	ecPrivKey, err := x509.ParseECPrivateKey(m5KeyParsedPEM.Bytes)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%#v", ecPrivKey)

}

func TestVerifyM5Sign(t *testing.T) {
	m5PubPEM := `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEWCB6hjZmL00xCFLwuGq+hs7+TL+X
2TcTzQ3oam2a09P0+I98Ni43nj3hSsukWjb6ypb/rnaRYMueBsVaGmgeyg==
-----END PUBLIC KEY-----`

	m5PubParsedPEM, _ := pem.Decode([]byte(m5PubPEM))
	pubKey, err := x509.ParsePKIXPublicKey(m5PubParsedPEM.Bytes)
	if err != nil {
		t.Error(err)
	}
	ecPubKey := pubKey.(*ecdsa.PublicKey)
	fmt.Printf("PubKey:\n%#v\n", ecPubKey)
	sum := sha256.Sum256(m5PubParsedPEM.Bytes)
	fmt.Printf("Pubkey hash:% 02X", sum)

	hashHex := "28 1A E7 AD 03 E6 0E F5 BE 96 10 CA 8A 79 FB 26 60 13 2A EA D4 B1 2F 3F 18 03 4C B0 AD D1 08 2A"
	hashHex = strings.Replace(hashHex, " ", "", -1)
	hash, err := hex.DecodeString(hashHex)
	if err != nil {
		t.Fatal(err)
	}

	signHex := "30 45 02 21 00 C6 EB 42 FA BF 70 DB 2A 03 C0 B7 D2 E9 E5 72 D9 F2 75 BE 54 EA 29 20 C1 5D 5B 3D 47 A2 80 D7 D5 02 20 37 14 C6 10 4A 29 8F FF 47 CF 9F F6 64 ED 07 A2 18 A4 BD 40 AA E2 A1 24 4E 17 A3 06 65 22 BC 90"
	signHex = strings.Replace(signHex, " ", "", -1)
	sign, err := hex.DecodeString(signHex)
	if err != nil {
		t.Fatal(err)
	}

	type SigData struct {
		R big.Int
		S big.Int
	}
	var sigData []*big.Int
	_, err = asn1.Unmarshal(sign, &sigData)
	if err != nil {
		t.Fatal(err)
	}

	isOk := ecdsa.Verify(ecPubKey, hash, sigData[0], sigData[1])
	if isOk {
		t.Log("Ok")
	} else {
		t.Log("Fail")
	}

}
