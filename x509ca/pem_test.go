package main

import (
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	// 	const pubPEM = `
	// -----BEGIN PUBLIC KEY-----
	// MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAlRuRnThUjU8/prwYxbtyWPT9pURI3lbsKMiB6Fn/VHOKE13p4D8xgOCADpdRagdT6n4etr9atzDKUSvpMtR3
	// CP5noNc97WiNCggBjVWhs7szEe8ugyqF23XwpHQ6uV1LKH50m92MbOWfCtjU9p/xqhNpQQ1AZhqNy5Gevap5k8XzRmjSldNAFZMY7Yv3Gi+nyCwGwpVtBUwhuLzgNFK/
	// yDtw2WcWmUU7NuC8Q6MWvPebxVtCfVp/iQU6q60yyt6aGOBkhAX0LpKAEhKidixYnP9PNVBvxgu3XZ4P36gZV6+ummKdBVnc3NqwBLu5+CcdRdusmHPHd5pHf4/38Z3/
	// 6qU2a/fPvWzceVTEgZ47QjFMTCTmCwNt29cvi7zZeQzjtwQgn4ipN9NibRH/Ax/qTbIzHfrJ1xa2RteWSdFjwtxi9C20HUkjXSeI4YlzQMH0fPX6KCE7aVePTOnB69I/
	// a9/q96DiXZajwlpq3wFctrs1oXqBp5DVrCIj8hU2wNgB7LtQ1mCtsYz//heai0K9PhE4X6hiE0YmeAZjR0uHl8M/5aW9xCoJ72+12kKpWAa0SFRWLy6FejNYCYpkupVJ
	// yecLk/4L1W0l6jQQZnWErXZYe0PNFcmwGXy1Rep83kfBRNKRy5tvocalLlwXLdUkAIU+2GKjyT3iMuzZxxFxPFMCAwEAAQ==
	// -----END PUBLIC KEY-----`

	// 	const pubPEM = `
	// -----BEGIN PUBLIC KEY-----
	// MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDsqXnN+EsT8uORBwBggrtP3LDdaCW
	// rdBApbZx5YuzHSOAE7mEPTAIce0o6zKEjZWLV8SOMhuKD0g3NJX+m+geByVoFijHBrJ
	// t7cZ/b+qJzLZ+0g6IXLfaZGQOCRu8K4LFBp1/5XHHyyTUYdkoA9lxKVYGVgaWGAKCax
	// q2wubHwzwIDAQAB
	// -----END PUBLIC KEY-----`

	const pubPEM = `
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDVrJE5eCnIrtxArvtYuZHHB4rFvD9
JLnJYD9oked1T00xPjhn3Qi+bwrmMizpONGijqUpUdpJs4i6H6h9QTLlgaGpEuJ+XDi
QbAOa6E9oCXFaPrr4WJ7zJ6CH05jfzNkxD6Bl4M+W47p80CBgPJxDIJb7F6SOPYTJP9
u35jMkOoQIDAQAB
-----END PUBLIC KEY-----`

	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		panic("failed to parse PEM block containing the public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic("failed to parse DER encoded public key: " + err.Error())
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		fmt.Println("pub is of type RSA:", pub)
	case *dsa.PublicKey:
		fmt.Println("pub is of type DSA:", pub)
	case *ecdsa.PublicKey:
		fmt.Println("pub is of type ECDSA:", pub)
	default:
		panic("unknown type of public key")
	}
}
