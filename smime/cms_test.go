package smime

import (
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"github.com/InfiniteLoopSpace/go_S-MIME/pki"
	"testing"
)
import "github.com/InfiniteLoopSpace/go_S-MIME/smime"

var (
	root = pki.New(pki.IsCA, pki.Subject(pkix.Name{
		CommonName: "root.example.com",
	}))

	intermediate = root.Issue(pki.IsCA, pki.Subject(pkix.Name{
		CommonName: "intermediate.example.com",
	}))

	leaf = intermediate.Issue(pki.Subject(pkix.Name{
		CommonName: "leaf.example.com",
	}))

	keyPair = tls.Certificate{
		Certificate: [][]byte{leaf.Certificate.Raw, intermediate.Certificate.Raw, root.Certificate.Raw},
		PrivateKey:  leaf.PrivateKey.(crypto.PrivateKey),
	}
)

func TestName(t *testing.T) {
	mail := "From: Alice\nTo: Bob\n\nHello World!"
	SMIMESender, _ := smime.New()
	ciphertext, _ := SMIMESender.Encrypt([]byte(mail), []*x509.Certificate{leaf.Certificate})
	// Bob
	BobkeyPair := keyPair
	SMIMERecipient, _ := smime.New(BobkeyPair)
	plaintext, _ := SMIMERecipient.Decrypt(ciphertext)
	t.Logf("%s", plaintext)
}
