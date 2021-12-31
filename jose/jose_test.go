package jose

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io"
	"testing"
)
import "github.com/go-jose/go-jose/v3"

func TestJWEMulti2(t *testing.T) {
	// Generate a public/private key pair to use for this example.
	alicePk, _ := rsa.GenerateKey(rand.Reader, 2048)
	bobPk, _ := rsa.GenerateKey(rand.Reader, 2048)

	// Instantiate an encrypter using RSA-OAEP with AES128-GCM. An error would
	// indicate that the selected algorithm(s) are not currently supported.
	encrypter, _ := jose.NewMultiEncrypter(jose.A128GCM, []jose.Recipient{
		{Algorithm: jose.RSA_OAEP, Key: &alicePk.PublicKey},
		{Algorithm: jose.RSA_OAEP, Key: &bobPk.PublicKey},
	}, nil)

	// Encrypt a sample plaintext. Calling the encrypter returns an encrypted
	// JWE object, which can then be serialized for output afterwards. An error
	// would indicate a problem in an underlying cryptographic primitive.
	var plaintext = []byte("Lorem ipsum dolor sit amet")
	var aad = []byte("This is authenticated, but public data")
	object, _ := encrypter.EncryptWithAuthData(plaintext, aad)

	// Serialize the encrypted object using the full serialization format.
	// Alternatively you can also use the compact format here by calling
	// object.CompactSerialize() instead.
	serialized := object.FullSerialize()
	fmt.Printf("Serialized message: %s\n", serialized)
	//pretty.Printf("Serialized message: %s\n", serialized)

	// Parse the serialized, encrypted JWE object. An error would indicate that
	// the given input did not represent a valid message.
	object, _ = jose.ParseEncrypted(serialized)

	// Now we can decrypt and get back our original plaintext. An error here
	// would indicate that the message failed to decrypt, e.g. because the auth
	// tag was broken or the message was tampered with.
	idx, header, decrypted, _ := object.DecryptMulti(alicePk)
	fmt.Printf("Decrypted by Alice (#%d): %s, Headers:%+v\n", idx, string(decrypted), header)
	idx, header, decrypted, _ = object.DecryptMulti(bobPk)
	fmt.Printf("Decrypted by Bob (#%d): %s, Headers:%+v\n", idx, string(decrypted), header)
	fmt.Printf("AAD: %s\n", string(object.GetAuthData()))
}

func TestJWEMulti(t *testing.T) {
	// Generate a public/private key pair to use for this example.
	alicePk, _ := rsa.GenerateKey(rand.Reader, 2048)
	bobPk, _ := rsa.GenerateKey(rand.Reader, 2048)

	// Instantiate an encrypter using RSA-OAEP with AES128-GCM. An error would
	// indicate that the selected algorithm(s) are not currently supported.
	encrypter, _ := jose.NewMultiEncrypter(jose.A128GCM, []jose.Recipient{
		{Algorithm: jose.RSA_OAEP, Key: &alicePk.PublicKey},
		{Algorithm: jose.RSA_OAEP, Key: &bobPk.PublicKey},
	}, nil)

	// Encrypt a sample plaintext. Calling the encrypter returns an encrypted
	// JWE object, which can then be serialized for output afterwards. An error
	// would indicate a problem in an underlying cryptographic primitive.
	var plaintext = []byte("Lorem ipsum dolor sit amet")
	var aad = []byte("This is authenticated, but public data")
	object, _ := encrypter.EncryptWithAuthData(plaintext, aad)

	// Serialize the encrypted object using the full serialization format.
	// Alternatively you can also use the compact format here by calling
	// object.CompactSerialize() instead.
	serialized := object.FullSerialize()
	fmt.Printf("Serialized message: %s\n", serialized)
	//pretty.Printf("Serialized message: %s\n", serialized)

	// Parse the serialized, encrypted JWE object. An error would indicate that
	// the given input did not represent a valid message.
	object, _ = jose.ParseEncrypted(serialized)

	// Now we can decrypt and get back our original plaintext. An error here
	// would indicate that the message failed to decrypt, e.g. because the auth
	// tag was broken or the message was tampered with.
	idx, header, decrypted, _ := object.DecryptMulti(alicePk)
	fmt.Printf("Decrypted by Alice (#%d): %s, Headers:%+v\n", idx, string(decrypted), header)
	idx, header, decrypted, _ = object.DecryptMulti(bobPk)
	fmt.Printf("Decrypted by Bob (#%d): %s, Headers:%+v\n", idx, string(decrypted), header)
	fmt.Printf("AAD: %s\n", string(object.GetAuthData()))
}

func TestJWERSA(t *testing.T) {
	// Generate a public/private key pair to use for this example.
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	// Instantiate an encrypter using RSA-OAEP with AES128-GCM. An error would
	// indicate that the selected algorithm(s) are not currently supported.
	publicKey := &privateKey.PublicKey
	encrypter, err := jose.NewEncrypter(jose.A128GCM, jose.Recipient{Algorithm: jose.RSA_OAEP, Key: publicKey}, nil)
	if err != nil {
		panic(err)
	}

	// Encrypt a sample plaintext. Calling the encrypter returns an encrypted
	// JWE object, which can then be serialized for output afterwards. An error
	// would indicate a problem in an underlying cryptographic primitive.
	var plaintext = []byte("Lorem ipsum dolor sit amet")
	object, err := encrypter.Encrypt(plaintext)
	if err != nil {
		panic(err)
	}

	// Serialize the encrypted object using the full serialization format.
	// Alternatively you can also use the compact format here by calling
	// object.CompactSerialize() instead.
	serialized := object.FullSerialize()
	fmt.Printf("Serialized message: %s\n", serialized)

	// Parse the serialized, encrypted JWE object. An error would indicate that
	// the given input did not represent a valid message.
	object, err = jose.ParseEncrypted(serialized)
	if err != nil {
		panic(err)
	}

	object.GetAuthData()

	// Now we can decrypt and get back our original plaintext. An error here
	// would indicate that the message failed to decrypt, e.g. because the auth
	// tag was broken or the message was tampered with.
	decrypted, err := object.Decrypt(privateKey)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Decrypted: %s\n", string(decrypted))
}

func TestJWEECC(t *testing.T) {
	// Generate a public/private key pair to use for this example.
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	publicKey := &privateKey.PublicKey
	encrypter, err := jose.NewEncrypter(jose.A128GCM, jose.Recipient{Algorithm: jose.ECDH_ES, Key: publicKey}, nil)
	if err != nil {
		panic(err)
	}

	// Encrypt a sample plaintext. Calling the encrypter returns an encrypted
	// JWE object, which can then be serialized for output afterwards. An error
	// would indicate a problem in an underlying cryptographic primitive.
	var plaintext = []byte("Lorem ipsum dolor sit amet")
	object, err := encrypter.Encrypt(plaintext)
	if err != nil {
		panic(err)
	}

	// Serialize the encrypted object using the full serialization format.
	// Alternatively you can also use the compact format here by calling
	// object.CompactSerialize() instead.
	serialized := object.FullSerialize()
	fmt.Printf("Serialized message: %s\n", serialized)

	// Parse the serialized, encrypted JWE object. An error would indicate that
	// the given input did not represent a valid message.
	object, err = jose.ParseEncrypted(serialized)
	if err != nil {
		panic(err)
	}

	object.GetAuthData()

	// Now we can decrypt and get back our original plaintext. An error here
	// would indicate that the message failed to decrypt, e.g. because the auth
	// tag was broken or the message was tampered with.
	decrypted, err := object.Decrypt(privateKey)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Decrypted: %s\n", string(decrypted))
}

func TestJWEPBE(t *testing.T) {

	// Instantiate an encrypter using RSA-OAEP with AES128-GCM. An error would
	// indicate that the selected algorithm(s) are not currently supported.
	encrypter, err := jose.NewEncrypter(jose.A128GCM, jose.Recipient{Algorithm: jose.PBES2_HS256_A128KW, Key: "Passw0rd"}, nil)
	if err != nil {
		panic(err)
	}

	// Encrypt a sample plaintext. Calling the encrypter returns an encrypted
	// JWE object, which can then be serialized for output afterwards. An error
	// would indicate a problem in an underlying cryptographic primitive.
	var plaintext = []byte("Lorem ipsum dolor sit amet")
	object, err := encrypter.Encrypt(plaintext)
	if err != nil {
		panic(err)
	}

	// Serialize the encrypted object using the full serialization format.
	// Alternatively you can also use the compact format here by calling
	// object.CompactSerialize() instead.
	serialized := object.FullSerialize()
	fmt.Printf("Serialized message: %s\n", serialized)

	// Parse the serialized, encrypted JWE object. An error would indicate that
	// the given input did not represent a valid message.
	object, err = jose.ParseEncrypted(serialized)
	if err != nil {
		panic(err)
	}

	object.GetAuthData()

	// Now we can decrypt and get back our original plaintext. An error here
	// would indicate that the message failed to decrypt, e.g. because the auth
	// tag was broken or the message was tampered with.
	decrypted, err := object.Decrypt("Passw0rd")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Decrypted: %s\n", string(decrypted))
}

func TestJWSRSA(t *testing.T) {
	// Generate a public/private key pair to use for this example.
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	// Instantiate a signer using RSASSA-PSS (SHA512) with the given private key.
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.PS512, Key: privateKey}, nil)
	if err != nil {
		panic(err)
	}

	// Sign a sample payload. Calling the signer returns a protected JWS object,
	// which can then be serialized for output afterwards. An error would
	// indicate a problem in an underlying cryptographic primitive.
	var payload = []byte("Lorem ipsum dolor sit amet")
	object, err := signer.Sign(payload)
	if err != nil {
		panic(err)
	}

	// Serialize the encrypted object using the full serialization format.
	// Alternatively you can also use the compact format here by calling
	// object.CompactSerialize() instead.
	serialized := object.FullSerialize()

	// Parse the serialized, protected JWS object. An error would indicate that
	// the given input did not represent a valid message.
	object, err = jose.ParseSigned(serialized)
	if err != nil {
		panic(err)
	}

	// Now we can verify the signature on the payload. An error here would
	// indicate that the message failed to verify, e.g. because the signature was
	// broken or the message was tampered with.
	output, err := object.Verify(&privateKey.PublicKey)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(output))

}

func TestJWEMulti_PBE_ECC_AES(t *testing.T) {
	// Generate a public/private key pair to use for this example.
	password := "Passw0rd"
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	aes128 := make([]byte, 16)
	_, err = io.ReadFull(rand.Reader, aes128)
	if err != nil {
		panic(err)
	}

	publicKey := &privateKey.PublicKey

	// Instantiate an encrypter using RSA-OAEP with AES128-GCM. An error would
	// indicate that the selected algorithm(s) are not currently supported.
	encrypter, err := jose.NewMultiEncrypter(jose.A128GCM, []jose.Recipient{
		{Algorithm: jose.PBES2_HS256_A128KW, Key: password},
		{Algorithm: jose.ECDH_ES_A256KW, Key: publicKey},
		{Algorithm: jose.A128GCMKW, Key: aes128},
	}, nil)
	if err != nil {
		panic(err)
	}

	// Encrypt a sample plaintext. Calling the encrypter returns an encrypted
	// JWE object, which can then be serialized for output afterwards. An error
	// would indicate a problem in an underlying cryptographic primitive.
	var plaintext = []byte("Lorem ipsum dolor sit amet")
	var aad = []byte("This is authenticated, but public data")
	object, _ := encrypter.EncryptWithAuthData(plaintext, aad)

	// Serialize the encrypted object using the full serialization format.
	// Alternatively you can also use the compact format here by calling
	// object.CompactSerialize() instead.
	serialized := object.FullSerialize()
	fmt.Printf("Serialized message: %s\n", serialized)

	// Parse the serialized, encrypted JWE object. An error would indicate that
	// the given input did not represent a valid message.
	object, _ = jose.ParseEncrypted(serialized)

	// Now we can decrypt and get back our original plaintext. An error here
	// would indicate that the message failed to decrypt, e.g. because the auth
	// tag was broken or the message was tampered with.
	idx, header, decrypted, _ := object.DecryptMulti(password)
	fmt.Printf("Decrypted by password (#%d): %s, Headers:%+v\n", idx, string(decrypted), header)
	idx, header, decrypted, _ = object.DecryptMulti(privateKey)
	fmt.Printf("Decrypted by ECC private key (#%d): %s, Headers:%+v\n", idx, string(decrypted), header)
	idx, header, decrypted, _ = object.DecryptMulti(aes128)
	fmt.Printf("Decrypted by AES key (#%d): %s, Headers:%+v\n", idx, string(decrypted), header)
	fmt.Printf("AAD: %s\n", string(object.GetAuthData()))
}
