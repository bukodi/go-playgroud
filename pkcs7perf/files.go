package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strconv"
	"sync"
	"text/template"
	"time"

	"github.com/fullsailor/pkcs7"
)

var n = flag.Int("n", 10, "Number of test files")
var srcDir = flag.String("src", "/tmp/in", "Input directory")
var dstDir = flag.String("dst", "/tmp/out", "Output directory")
var saveFiles = flag.Bool("saveFiles", false, "Save the files or use only in memory")

var cert certKeyPair
var wg sync.WaitGroup

func main() {
	start := time.Now()
	flag.Parse()
	flag.Usage()

	fmt.Println("src:" + *srcDir)
	fmt.Println("dst:" + *dstDir)

	cert, _ = createTestCertificate()

	wg.Add(*n)
	for i := 0; i < *n; i++ {
		go createAndProcessFile(i)
	}
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("Done. Elapsed time: %s \n", elapsed)
}

func createAndProcessFile(index int) {
	defer wg.Done()
	data := createTestFile(index)
	signature := signData(data)
	if *saveFiles {
		err := ioutil.WriteFile(*srcDir+"/data"+strconv.Itoa(index)+".xml", data, 0644)
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile(*srcDir+"/data"+strconv.Itoa(index)+".sig", signature, 0644)
		if err != nil {
			panic(err)
		}
	}

}

func signData(data []byte) (signature []byte) {
	// Initialize a SignedData struct with content to be signed
	signedData, err := pkcs7.NewSignedData(data)
	if err != nil {
		fmt.Printf("Cannot initialize signed data: %s", err)
	}

	// Add the signing cert and private key
	if err := signedData.AddSigner(cert.Certificate, cert.PrivateKey, pkcs7.SignerInfoConfig{}); err != nil {
		fmt.Printf("Cannot add signer: %s", err)
	}

	// Call Detach() is you want to remove content from the signature
	// and generate an S/MIME detached signature
	signedData.Detach()

	// Finish() to obtain the signature bytes
	detachedSignature, err := signedData.Finish()
	if err != nil {
		fmt.Printf("Cannot finish signing data: %s", err)
	}
	var sigBuff bytes.Buffer
	pem.Encode(&sigBuff, &pem.Block{Type: "PKCS7", Bytes: detachedSignature})
	return sigBuff.Bytes()
}

type certKeyPair struct {
	Certificate *x509.Certificate
	PrivateKey  *rsa.PrivateKey
}

func createTestCertificate() (certKeyPair, error) {
	signer, err := createTestCertificateByIssuer("Eddard Stark", nil)
	if err != nil {
		return certKeyPair{}, err
	}
	fmt.Println("Created root cert")
	pem.Encode(os.Stdout, &pem.Block{Type: "CERTIFICATE", Bytes: signer.Certificate.Raw})
	pair, err := createTestCertificateByIssuer("Jon Snow", signer)
	if err != nil {
		return certKeyPair{}, err
	}
	fmt.Println("Created signer cert")
	pem.Encode(os.Stdout, &pem.Block{Type: "CERTIFICATE", Bytes: pair.Certificate.Raw})
	return *pair, nil
}

func createTestCertificateByIssuer(name string, issuer *certKeyPair) (*certKeyPair, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return nil, err
	}
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 32)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, err
	}

	template := x509.Certificate{
		SerialNumber:       serialNumber,
		SignatureAlgorithm: x509.SHA256WithRSA,
		Subject: pkix.Name{
			CommonName:   name,
			Organization: []string{"Acme Co"},
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(1, 0, 0),
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageEmailProtection},
	}
	var issuerCert *x509.Certificate
	var issuerKey crypto.PrivateKey
	if issuer != nil {
		issuerCert = issuer.Certificate
		issuerKey = issuer.PrivateKey
	} else {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
		issuerCert = &template
		issuerKey = priv
	}
	cert, err := x509.CreateCertificate(rand.Reader, &template, issuerCert, priv.Public(), issuerKey)
	if err != nil {
		return nil, err
	}
	leaf, err := x509.ParseCertificate(cert)
	if err != nil {
		return nil, err
	}
	return &certKeyPair{
		Certificate: leaf,
		PrivateKey:  priv,
	}, nil
}

type testFileParams struct {
	Index int
}

func createTestFile(index int) []byte {
	params := testFileParams{index}
	var buff bytes.Buffer
	fileTemplate.Execute(&buff, params)
	return buff.Bytes()
}

var fileTemplate, _ = template.New("pacs008").Parse(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<ns2:Document xmlns:ns2="urn:iso:std:iso:20022:tech:xsd:pacs.008.001.02">
	<FIToFICstmrCdtTrf>
		<GrpHdr>
			<MsgId>{{.Index}}</MsgId>
			<CreDtTm>2017-11-07T18:09:59.904Z</CreDtTm>
			<NbOfTxs>1</NbOfTxs>
			<TtlIntrBkSttlmAmt Ccy="HUF">3376033</TtlIntrBkSttlmAmt>
			<IntrBkSttlmDt>2017-11-07</IntrBkSttlmDt>
			<SttlmInf>
				<SttlmMtd>CLRG</SttlmMtd>
			</SttlmInf>
			<PmtTpInf>
				<SvcLvl>
					<Cd>VIPS</Cd>
				</SvcLvl>
				<LclInstrm>
					<Cd>INST</Cd>
				</LclInstrm>
			</PmtTpInf>
		</GrpHdr>
		<CdtTrfTxInf>
			<PmtId>
				<InstrId>d6b3fe80-e98b-47a3-9e2a-09c0d931cd27</InstrId>
				<EndToEndId>d6b3fe80-e98b-47a3-9e2a-09c0d931cd27</EndToEndId>
				<TxId>d6b3fe80-e98b-47a3-9e2a-09c0d931cd27</TxId>
			</PmtId>
			<IntrBkSttlmAmt Ccy="HUF">3376033</IntrBkSttlmAmt>
			<AccptncDtTm>2017-11-07T18:09:59.904Z</AccptncDtTm>
			<ChrgBr>SLEV</ChrgBr>
			<Dbtr>
				<Nm>Fel-Adó Andrea</Nm>
			</Dbtr>
			<DbtrAcct>
				<Nm>HU361234567890</Nm>
			</DbtrAcct>
			<DbtrAgt>
				<FinInstnId>
					<BIC>XYZZHUHB</BIC>
				</FinInstnId>
			</DbtrAgt>
			<CdtrAgt>
				<FinInstnId>
					<BIC>ABCDHUHB</BIC>
				</FinInstnId>
			</CdtrAgt>
			<Cdtr>
				<Nm>Megkapó Anikó</Nm>
			</Cdtr>
			<CdtrAcct>
				<Nm>HU360987654321</Nm>
			</CdtrAcct>
			<InstrForNxtAgt>
				<InstrInf>/STTLMREF/VCSM17117000001</InstrInf>
			</InstrForNxtAgt>
		</CdtTrfTxInf>
	</FIToFICstmrCdtTrf>
</ns2:Document>`)
