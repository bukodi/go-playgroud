package main

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net/http"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

var (
	mySSLKey crypto.Signer
)

const cert_fakelerootx1 = `-----BEGIN CERTIFICATE-----
MIIFATCCAumgAwIBAgIRAKc9ZKBASymy5TLOEp57N98wDQYJKoZIhvcNAQELBQAw
GjEYMBYGA1UEAwwPRmFrZSBMRSBSb290IFgxMB4XDTE2MDMyMzIyNTM0NloXDTM2
MDMyMzIyNTM0NlowGjEYMBYGA1UEAwwPRmFrZSBMRSBSb290IFgxMIICIjANBgkq
hkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA+pYHvQw5iU3v2b3iNuYNKYgsWD6KU7aJ
diddtZQxSWYzUI3U0I1UsRPTxnhTifs/M9NW4ZlV13ZfB7APwC8oqKOIiwo7IwlP
xg0VKgyz+kT8RJfYr66PPIYP0fpTeu42LpMJ+CKo9sbpgVNDZN2z/qiXrRNX/VtG
TkPV7a44fZ5bHHVruAxvDnylpQxJobtCBWlJSsbIRGFHMc2z88eUz9NmIOWUKGGj
EmP76x8OfRHpIpuxRSCjn0+i9+hR2siIOpcMOGd+40uVJxbRRP5ZXnUFa2fF5FWd
O0u0RPI8HON0ovhrwPJY+4eWKkQzyC611oLPYGQ4EbifRsTsCxUZqyUuStGyp8oa
aoSKfF6X0+KzGgwwnrjRTUpIl19A92KR0Noo6h622OX+4sZiO/JQdkuX5w/HupK0
A0M0WSMCvU6GOhjGotmh2VTEJwHHY4+TUk0iQYRtv1crONklyZoAQPD76hCrC8Cr
IbgsZLfTMC8TWUoMbyUDgvgYkHKMoPm0VGVVuwpRKJxv7+2wXO+pivrrUl2Q9fPe
Kk055nJLMV9yPUdig8othUKrRfSxli946AEV1eEOhxddfEwBE3Lt2xn0hhiIedbb
Ftf/5kEWFZkXyUmMJK8Ra76Kus2ABueUVEcZ48hrRr1Hf1N9n59VbTUaXgeiZA50
qXf2bymE6F8CAwEAAaNCMEAwDgYDVR0PAQH/BAQDAgEGMA8GA1UdEwEB/wQFMAMB
Af8wHQYDVR0OBBYEFMEmdKSKRKDm+iAo2FwjmkWIGHngMA0GCSqGSIb3DQEBCwUA
A4ICAQBCPw74M9X/Xx04K1VAES3ypgQYH5bf9FXVDrwhRFSVckria/7dMzoF5wln
uq9NGsjkkkDg17AohcQdr8alH4LvPdxpKr3BjpvEcmbqF8xH+MbbeUEnmbSfLI8H
sefuhXF9AF/9iYvpVNC8FmJ0OhiVv13VgMQw0CRKkbtjZBf8xaEhq/YqxWVsgOjm
dm5CAQ2X0aX7502x8wYRgMnZhA5goC1zVWBVAi8yhhmlhhoDUfg17cXkmaJC5pDd
oenZ9NVhW8eDb03MFCrWNvIh89DDeCGWuWfDltDq0n3owyL0IeSn7RfpSclpxVmV
/53jkYjwIgxIG7Gsv0LKMbsf6QdBcTjhvfZyMIpBRkTe3zuHd2feKzY9lEkbRvRQ
zbh4Ps5YBnG6CKJPTbe2hfi3nhnw/MyEmF3zb0hzvLWNrR9XW3ibb2oL3424XOwc
VjrTSCLzO9Rv6s5wi03qoWvKAQQAElqTYRHhynJ3w6wuvKYF5zcZF3MDnrVGLbh1
Q9ePRFBCiXOQ6wPLoUhrrbZ8LpFUFYDXHMtYM7P9sc9IAWoONXREJaO08zgFtMp4
8iyIYUyQAbsvx8oD2M8kRvrIRSrRJSl6L957b4AFiLIQ/GgV2curs0jje7Edx34c
idWw1VrejtwclobqNMVtG3EiPUIpJGpbMcJgbiLSmKkrvQtGng==
-----END CERTIFICATE-----`

const privateKeyPem = `-----BEGIN PRIVATE KEY-----
MIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQgcZXcRMQkNaqaaTWy
M1+LKl5brVVXRQSQS9dgyrmC5xygCgYIKoZIzj0DAQehRANCAASBYGGy1WmHSn/l
bHIosuRaApUXhCYYlZwZZVpgvKyZdAtEqd+Hd0oIj1QNwV2XDijOiSvyBhOLSn8Q
jBNKc/FY
-----END PRIVATE KEY-----`

type P11Manager struct {
	autocert.Manager
	TLSPrivateKey crypto.Signer
}

type P11EcdsaPrivKey struct {
	publicKey ecdsa.PublicKey
	originalD *big.Int
}

func (p11Key *P11EcdsaPrivKey) Public() crypto.PublicKey {
	return &p11Key.publicKey
}

func (p11Key *P11EcdsaPrivKey) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	normalKey := ecdsa.PrivateKey{
		PublicKey: p11Key.publicKey,
		D:         p11Key.originalD,
	}
	return normalKey.Sign(rand, digest, opts)
}

func NewP11EcdsaPrivKey(pkcs8Pem string) (*P11EcdsaPrivKey, error) {
	privateKeyBlock, _ := pem.Decode([]byte(privateKeyPem))
	key, _ := x509.ParsePKCS8PrivateKey(privateKeyBlock.Bytes)
	importedKey, _ := key.(*ecdsa.PrivateKey)

	p11Key := P11EcdsaPrivKey{
		publicKey: importedKey.PublicKey,
		originalD: importedKey.D,
	}
	return &p11Key, nil
}

func main() {
	//privateKeyBlock, _ := pem.Decode([]byte(privateKeyPem))
	//key, _ := x509.ParsePKCS8PrivateKey(privateKeyBlock.Bytes)
	//	importedKey, _ := key.(*ecdsa.PrivateKey)

	mySSLKey, _ = NewP11EcdsaPrivKey(privateKeyPem)
	//mySSLKey, _ = rsa.GenerateKey(rand.Reader, 2048)
	//mySSLKey, _ = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	fmt.Printf("Selected key:\n %#v\n\n", mySSLKey)
	//x509.MarshalPKCS1PrivateKey(*rsa.PrivateKey(&mySSLKey))

	acmeClient := acme.Client{
		//DirectoryURL: "https://acme.api.letsencrypt.org/directory",
		DirectoryURL: "https://acme-staging.api.letsencrypt.org/directory",
		//DirectoryURL: "https://acme-staging-v02.api.letsencrypt.org/directory",
		//Key: mySSLKey,
	}

	certManager := autocert.Manager{
		Client:     &acmeClient,
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("bukodi.ddns.net"), //Your domain here
		//Cache:      autocert.DirCache("/tmp/certs6"),          //Folder for storing certificates
		TLSKey: mySSLKey,
	}

	// certManager := P11Manager{
	// 	Manager:       baseCertManager,
	// 	TLSPrivateKey: mySSLKey,
	// }

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	})

	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	rootCAs.AppendCertsFromPEM([]byte(cert_fakelerootx1))

	certMgrTLSConfig := certManager.TLSConfig()

	wrappedTLSConfig := tls.Config{
		GetCertificate: certMgrTLSConfig.GetCertificate,
		NextProtos:     certMgrTLSConfig.NextProtos,
		RootCAs:        rootCAs,
	}

	server := &http.Server{
		Addr:      ":10443",
		TLSConfig: &wrappedTLSConfig,
	}

	go http.ListenAndServe("https://localhost:10443", certManager.HTTPHandler(nil))
	//http.HandleFunc("/", helloHandler)
	server.ListenAndServeTLS("", "") //Key and cert are coming from Let's Encrypt
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
}
