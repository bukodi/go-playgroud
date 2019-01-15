// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package winkerb2

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/alexbrainman/sspi/negotiate"
)

var (
	testURL = "https://noreg-nb33.noreg.intra/debug" // flag.String("url", "", "server URL for TestNegotiateHTTPClient")
)

// TODO: perhaps add Transport that is similar to http.Transport
// TODO: perhaps implement separate NTLMTransport and KerberosTransport (not sure about this idea)
// TODO: KerberosTransport is (I beleive) sinlge leg protocol, so it can be implemented easily (unlike NTLM)
// TODO: perhaps implement both server and client Transport

type httpClient struct {
	client    *http.Client
	transport *http.Transport
	url       string
}

const nb33CertPem = `-----BEGIN CERTIFICATE-----
MIIBuTCCASKgAwIBAgIEW9dB3jANBgkqhkiG9w0BAQsFADAhMR8wHQYDVQQDDBZu
b3JlZy1uYjMzLm5vcmVnLmludHJhMB4XDTE4MTAyOTE3MjMxMloXDTI4MTAyOTE3
MjMxMlowITEfMB0GA1UEAwwWbm9yZWctbmIzMy5ub3JlZy5pbnRyYTCBnzANBgkq
hkiG9w0BAQEFAAOBjQAwgYkCgYEAt4v+OEtW9EibtX7sgxae/Bf6w/NZ9pUBlzBG
XDFHKYSlEdN4NR5Pw/HdSdB8qco9CgJm27b0+kFX5xaPPkjY3/+rTFuLrMkQ+yXo
XIaMlYUCAK10JkUF5+moxjmhdTSDEayOryCk5JnCKU0z6xjqQ8svLlPgCkoiav/+
OXzX+rECAwEAATANBgkqhkiG9w0BAQsFAAOBgQAeY1EAORWiZIImrgrj2xm9AS6E
gwkeNaGKJ1udmLRVWsStWGAt9zTCSHIBWfW6Wnerv2C3inCBhIRvGUQ2H2w3HXQB
5I+9u2I/xU7OBR3kJwwrnhyGDvrXHKBbIOr3ctha+HqonMyYiuJGj+LfZ2ew1421
cKwpKl2Twyj5F7duyQ==
-----END CERTIFICATE-----`

func newHTTPClient(url string) *httpClient {

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM([]byte(nb33CertPem))

	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSClientConfig: &tls.Config{RootCAs: certPool},
	}
	return &httpClient{
		client:    &http.Client{Transport: transport},
		transport: transport,
		url:       url,
	}
}

func (c *httpClient) CloseIdleConnections() {
	c.transport.CloseIdleConnections()
}

func (c *httpClient) get(req *http.Request) (*http.Response, string, error) {
	res, err := c.client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, "", err
	}
	return res, string(body), nil
}

func (c *httpClient) canDoNegotiate() error {
	req, err := http.NewRequest("GET", c.url, nil)
	if err != nil {
		return err
	}
	res, _, err := c.get(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusUnauthorized {
		return fmt.Errorf("Unauthorized expected, but got %v", res.StatusCode)
	}
	authHeaders, found := res.Header["Www-Authenticate"]
	if !found {
		return fmt.Errorf("Www-Authenticate not found")
	}
	for _, h := range authHeaders {
		if h == "Negotiate" {
			return nil
		}
	}
	return fmt.Errorf("Www-Authenticate header does not contain Negotiate, but has %v", authHeaders)
}

func findAuthHeader(res *http.Response) ([]byte, error) {
	authHeaders, found := res.Header["Www-Authenticate"]
	if !found {
		return nil, fmt.Errorf("Www-Authenticate not found")
	}
	if len(authHeaders) != 1 {
		return nil, fmt.Errorf("Only one Www-Authenticate header expected, but %d found: %v", len(authHeaders), authHeaders)
	}
	if len(authHeaders[0]) < 10 {
		return nil, fmt.Errorf("Www-Authenticate header is to short: %q", authHeaders[0])
	}
	if !strings.HasPrefix(authHeaders[0], "Negotiate ") {
		return nil, fmt.Errorf("Www-Authenticate header is suppose to starts with \"Negotiate \", but is %q", authHeaders[0])
	}
	token, err := base64.StdEncoding.DecodeString(authHeaders[0][10:])
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (c *httpClient) startAuthorization(inputToken []byte) ([]byte, error) {
	req, err := http.NewRequest("GET", c.url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Negotiate "+base64.StdEncoding.EncodeToString(inputToken))
	res, _, err := c.get(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusUnauthorized {
		return nil, fmt.Errorf("Unauthorized expected, but got %v", res.StatusCode)
	}
	outputToken, err := findAuthHeader(res)
	if err != nil {
		return nil, err
	}
	return outputToken, nil
}

func (c *httpClient) completeAuthorization(inputToken []byte) (*http.Response, string, error) {
	req, err := http.NewRequest("GET", c.url, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Authorization", "Negotiate "+base64.StdEncoding.EncodeToString(inputToken))
	res, body, err := c.get(req)
	if err != nil {
		return nil, "", err
	}
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusUnauthorized {
		return nil, "", fmt.Errorf("OK expected, but got %v", res.StatusCode)
	}
	return res, body, nil
}

func TestKerberosHTTPClient(t *testing.T) {
	// TODO: combine client and server tests so we don't need external server
	if len(testURL) == 0 {
		t.Skip("Skipping due to empty \"url\" parameter")
	}

	u, err := url.Parse(testURL)
	if err != nil {
		t.Fatal(err)
	}
	targetName := "http/" + strings.ToUpper(u.Host)

	cred, err := negotiate.AcquireCurrentUserCredentials()
	if err != nil {
		t.Fatal(err)
	}
	defer cred.Release()

	secctx, token, err := negotiate.NewClientContext(cred, targetName)
	if err != nil {
		t.Fatal(err)
	}
	defer secctx.Release()

	client := newHTTPClient(testURL)
	defer client.CloseIdleConnections()

	err = client.canDoNegotiate()
	if err != nil {
		t.Fatal(err)
	}
	res, _, err := client.completeAuthorization(token)
	if err != nil {
		t.Fatal(err)
	}
	serverToken, err := findAuthHeader(res)
	if err != nil {
		t.Fatal(err)
	}
	authCompleted, lastToken, err := secctx.Update(serverToken)
	if err != nil {
		t.Fatal(err)
	}
	if !authCompleted {
		t.Fatal("client authentication should be completed now")
	}
	if len(lastToken) > 0 {
		t.Fatalf("last token supposed to be empty, but %v returned", lastToken)
	}
}

// TODO: See http://www.innovation.ch/personal/ronald/ntlm.html#connections about needed to keep connection alive during authentication.

func TestNegotiateHTTPServer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: implement Negotiate authentication here
		w.Write([]byte("hello"))
	}))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "hello" {
		t.Errorf("got %q, want hello", string(got))
	}
}
