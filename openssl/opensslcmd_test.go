package openssl

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
	"testing"
)

func TestOutputPipe(t *testing.T) {

	// Create named pipe
	keyOutPipe := "/tmp/0814/test.key"
	syscall.Mkfifo(keyOutPipe, 0600)
	defer os.Remove(keyOutPipe)

	go func() {
		cmd := exec.Command("openssl", "genrsa", "-out", keyOutPipe, "2048")
		// Just to forward the stdout
		cmd.Stdout = os.Stdout
		cmd.Run()
	}()

	// Open named pipe for reading
	fmt.Println("Opening named pipe for reading")
	keyOutFile, _ := os.OpenFile(keyOutPipe, os.O_RDONLY, 0600)
	fmt.Println("Reading")

	var buff bytes.Buffer
	fmt.Println("Waiting for someone to write something")
	io.Copy(&buff, keyOutFile)
	keyOutFile.Close()
	fmt.Printf("Data: %s\n", buff.String())
}

func TestInputPipe(t *testing.T) {

	// Create named pipe
	keyInPipe := "/tmp/0814/testIn.key"
	certOutPipe := "/tmp/0814/certOut.pem"
	syscall.Mkfifo(keyInPipe, 0600)
	syscall.Mkfifo(certOutPipe, 0600)
	defer os.Remove(keyInPipe)
	defer os.Remove(certOutPipe)

	go func() {
		cmd := exec.Command("openssl", "req",
			"-x509", "-new", "-nodes",
			"-key", keyInPipe,
			"-sha256", "-days", "1024",
			"-out", certOutPipe,
			"-subj", "/C=GB/CN=foo")
		// Just to forward the stdout
		cmd.Stdout = os.Stdout
		cmd.Run()
	}()

	// Open named pipe for reading
	fmt.Println("Writing key to openssl")
	keyInBuff := bytes.NewBufferString(testKey)
	keyInFile, _ := os.OpenFile(keyInPipe, os.O_WRONLY, 0600)
	io.Copy(keyInFile, keyInBuff)
	keyInFile.Close()
	fmt.Println("Key wrote to openssl")

	// Open named pipe for reading
	fmt.Println("Reading generated cert form openssl")
	certOutFile, _ := os.OpenFile(certOutPipe, os.O_RDONLY, 0600)
	var certOutBuff bytes.Buffer
	io.Copy(&certOutBuff, certOutFile)
	certOutFile.Close()
	fmt.Printf("Generated cert:\n%s\n", certOutBuff.String())
}

const testKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAuxBHEw3+v8Xm/9rf5UgCSTNctzLIG1Bjm1mms9OYLqczKS73
l9WFpB8F4y7EgOMRnZYKp6E7yuMHvvSHwya85wPoCGS8g5eZ0/RRHQMMn9tYyTI5
LFLbHH42FcsB39RiGjK2uXOc6OMTIT+rUqqsUXsdJLoA8VXPTiEUNVYYjtnL2xO0
5BonNgu1UlFRZSvecphE9hwc4NR5QQylUhtazxy/WPVbSCvPf/3Wns1XdMD7XnWR
fo2Otr2inTKnXMantwj4+Fm1+5R5uBf0oXLiNNQcnMMHmhg7+WWmKsBriKFoZC3T
MpAbw06o3uWAQ0Q6sC/0W/6XD97fA4sBwW/5qQIDAQABAoIBAAai5DKb237IMZLA
HBNRQ6t/I/nn1kuJxY7cVlqo1gxJqDn8zZHYZF5XL2lI3nXIGHbjvMsHoExpU3wF
xs84j5kOfWvWzw1IEo//aeVCl28QZAz3OCoHYniXTanmQtHDAhv10p+vp1BnxeT3
EkfjgCt/15/W7XOiXLFj4QinXkWrCsnaRfu+Wb0vhiUEIN81R9pgvRm+kHbzgu5W
eZnR1+gf2FxiZr+zWg9UNmVKy8mLS7RRAoCwufA5JTBF08XuiJ+4qnUCQookUeZ3
l+0PzP+nK7Xmxr6hAZ1Y59OMVVjeypk81U6m9R6zQfJauPppvdsmz35PmG2leBgG
7lm1V9ECgYEA7PpsRHqdzP3cvmifOEOXYoxJCUcGWN5i7HlStZkg0Quk12oU6USo
uvA63DnytWtzz7XMMp5nh9XcnlJRBT3ZLQD+YXFkH1V+9VC1i2XyJUEOTAXMU6fl
7XNhH5gCkNh12XDkeX3ZMidSiehfxuABXiQVdKnRSZapPEwrGFgaDM0CgYEAyhQt
dUFCfipkHC4K228UM7TqisxM6BFruik8pAvXOUdLXn5ffqo0KNtD8cAAgE5UaQgL
059hTjvRjm5xwkx6d0UKBY94AWio1iCH6hig/jIvXi4Wc+3R02JlSoO/2QOdzGMb
tLYd417xlumtGyX2TqFkFvhqcObHHcMLFmACoE0CgYB5RYYWXTFX8CoA/wVMA7r6
4ZOWvdQPsm6pWUTsTdqvX+gRnOXqogo+8CUPAlCkasKvbvd6h/mvV9A47SMtLYNw
Nmv3bdGw/02jOJRPK/KJAgvQ976iqO9PXpY7Vs0pVryoc89YJQD7W4gvrs0ktwm8
JXcdZrIFmKYuh0Qehyd9mQKBgFRJ5kwqVFnbxLYcXlr5EiwfIlWSseF6orybxreG
WNeDbWSUwbBLvkXsb4K+23apNXw55vT2XdgMC3SljL3GuK5XFb8MALpVtVbbatWy
QDTHKgrWnnbsk8DgIe/a1ILoh0FhdYUDEaRtTcfs4E+angpeNyl9pKhDGnrHiDBl
C7NhAoGASAwNtLkD6/cuBWBsTbjrjBlxPr61ZBrwIT88APNKB6lCatJV9VncfCIt
lYPWoAim72i7JZlhzW6U3n1tXH8nKCvHJiGsBwObVsYqXDtb6t8mknD0qelOxcJS
emoj0YbJPdXV8T1Qidgsjlyyiq7MDdtkJM89ngEQjZS5tAS8wTs=
-----END RSA PRIVATE KEY-----`
