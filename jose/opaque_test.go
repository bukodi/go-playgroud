package jose

import (
	"bufio"
	"bytes"
	"context"
	"crypto/aes"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"gopkg.in/square/go-jose.v2"
	josecipher "gopkg.in/square/go-jose.v2/cipher"
	"hash"
	"sync"
	"testing"
)

type pbkdf2CacheEntry struct {
	derivedKey []byte
	pswCtx     context.Context
}

func pbkdf2ParamsHash(password, salt []byte, iter, keyLen int, h func() hash.Hash) string {
	var buff bytes.Buffer
	w := bufio.NewWriter(&buff)
	w.Write(password)
	w.Write(salt)
	binary.Write(w, binary.LittleEndian, iter)
	binary.Write(w, binary.LittleEndian, keyLen)
	summ := h().Sum(buff.Bytes())
	return base64.RawStdEncoding.EncodeToString(summ)
}

var pbkdf2Cache = make(map[string]*pbkdf2CacheEntry)
var pbkdf2CacheLock sync.RWMutex

type myCachingPBEDecrypter struct {
	pswCallback PswCallback
}

func readP2S(p2sIf interface{}) ([]byte, error) {
	if p2sIf == nil {
		return nil, fmt.Errorf("square/go-jose: invalid P2S: must be present")
	}
	p2s := fmt.Sprintf("%s", p2sIf)
	p2sBytes, err := base64.RawURLEncoding.DecodeString(p2s)
	if err != nil {
		return nil, fmt.Errorf("square/go-jose: invalid P2S: invalid base64.: %w", err)
	}
	return p2sBytes, nil
}

func (cpd *myCachingPBEDecrypter) DecryptKey(encryptedKey []byte, header jose.Header) ([]byte, error) {
	alg := jose.KeyAlgorithm(header.Algorithm)
	if !(alg == jose.PBES2_HS256_A128KW || alg == jose.PBES2_HS512_A256KW || alg == jose.PBES2_HS384_A192KW) {
		return nil, jose.ErrUnsupportedAlgorithm
	}

	p2sBytes, err := readP2S(header.ExtraHeaders["p2s"])
	if err != nil {
		return nil, fmt.Errorf("square/go-jose: invalid P2S: invalid base64.: %w", err)
	}
	_ = p2sBytes
	p2cIf := header.ExtraHeaders["p2c"]
	if p2cIf == nil {

	}
	var p2c int
	switch i := p2cIf.(type) {
	case nil:
		return nil, fmt.Errorf("square/go-jose: invalid P2C: must be present")
	case float64:
		p2c = int(i)
	case float32:
		p2c = int(i)
	case int64:
		p2c = int(i)
	default:
		return nil, fmt.Errorf("square/go-jose: invalid P2C: must be present")
	}
	if p2c <= 0 {
		return nil, fmt.Errorf("square/go-jose: invalid P2C: must be a positive integer")
	}

	// salt is UTF8(Alg) || 0x00 || Salt Input
	salt := bytes.Join([][]byte{[]byte(alg), p2sBytes}, []byte{0x00})

	// derive key
	keyLen, h := getPbkdf2Params(alg)
	password := cpd.pswCallback("")
	pbkdf2ParamsHash := pbkdf2ParamsHash(password, salt, p2c, keyLen, h)
	pbkdf2CacheLock.RLock()
	ce := pbkdf2Cache[pbkdf2ParamsHash]
	pbkdf2CacheLock.RUnlock()
	if ce == nil {
		key := pbkdf2.Key(password, salt, p2c, keyLen, h)
		ce = &pbkdf2CacheEntry{key, nil}
		pbkdf2CacheLock.Lock()
		pbkdf2Cache[pbkdf2ParamsHash] = ce
		pbkdf2CacheLock.Unlock()
	}

	// use AES cipher with derived key
	block, err := aes.NewCipher(ce.derivedKey)
	if err != nil {
		return nil, err
	}

	cek, err := josecipher.KeyUnwrap(block, encryptedKey)
	if err != nil {
		return nil, err
	}
	return cek, nil
}

// getPbkdf2Params returns the key length and hash function used in
// pbkdf2.Key.
func getPbkdf2Params(alg jose.KeyAlgorithm) (int, func() hash.Hash) {
	switch alg {
	case jose.PBES2_HS256_A128KW:
		return 16, sha256.New
	case jose.PBES2_HS384_A192KW:
		return 24, sha512.New384
	case jose.PBES2_HS512_A256KW:
		return 32, sha512.New
	default:
		panic("invalid algorithm")
	}
}

func TestJWEPBEForOpaque(t *testing.T) {

	encrypter, err := jose.NewEncrypter(jose.A128GCM, jose.Recipient{Algorithm: jose.PBES2_HS256_A128KW, Key: "Passw0rd"}, nil)
	if err != nil {
		panic(err)
	}

	var plaintext = []byte("Lorem ipsum dolor sit amet")
	object, err := encrypter.Encrypt(plaintext)
	if err != nil {
		panic(err)
	}

	serialized := object.FullSerialize()
	fmt.Printf("Serialized message: %q\n", serialized)

	object, err = jose.ParseEncrypted(serialized)
	if err != nil {
		panic(err)
	}

	object.GetAuthData()

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		decrypted1, err := object.Decrypt("Passw0rd")
		if err != nil {
			t.Error(err)
		} else {
			t.Logf("Decrypted1: %s", string(decrypted1))
		}
		wg.Done()
	}()

	go func() {
		pswCallBack := NewPswCallback("Passw0rd")
		mcpd := &myCachingPBEDecrypter{pswCallBack}
		decrypted2, err := object.Decrypt(mcpd)
		if err != nil {
			t.Error(err)
		} else {
			t.Logf("Decrypted2: %s", string(decrypted2))
		}
		wg.Done()
	}()
	wg.Wait()
}

const testJWE = "{\"protected\":\"eyJhbGciOiJQQkVTMi1IUzI1NitBMTI4S1ciLCJlbmMiOiJBMTI4R0NNIiwicDJjIjoxMDAwMDAsInAycyI6IlFCSXdEX21DaVpTUFgtcndFclpQY2cifQ\",\"encrypted_key\":\"WBgaIeo37t5MHlToLBQPqFcbFxVwSCbB\",\"iv\":\"JvMtYv-A82A9MPd1\",\"ciphertext\":\"7xpWzNkzI2dXW1XPJ4x5rboH2Ap30k7VAbk\",\"tag\":\"A6Dw_EqJmUKkq84IMZcQiw\"}"

func TestPDEDecrypter(t *testing.T) {
	object, _ := jose.ParseEncrypted(testJWE)
	decrypted1, err := object.Decrypt("Passw0rd")
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("Decrypted1: %s", string(decrypted1))
	}

}

func TestOpaqueDecrypter(t *testing.T) {
	object, _ := jose.ParseEncrypted(testJWE)
	pswCallBack := NewPswCallback("Passw0rd")
	mcpd := &myCachingPBEDecrypter{pswCallBack}
	decrypted2, err := object.Decrypt(mcpd)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("Decrypted2: %s", string(decrypted2))
	}
}
