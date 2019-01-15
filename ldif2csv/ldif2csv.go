package main

import (
	"crypto/x509"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-ldap/ldif"
)

var (
	csvHeaders = []string{"pos", "mail", "dn", "Issuer DN", "Subject DN", "Serial", "Valid from", "Valid to"}
)

func main() {
	if len(os.Args) != 3 {

		log.Fatal("Usage: " + filepath.Base(os.Args[0]) + " <inputFile> <ouputFile>")
	}

	csvFile, err := os.Create(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()
	csvWriter.Write(csvHeaders)

	ldifbytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	ldifStr := string(ldifbytes)

	l, err := ldif.Parse(ldifStr)
	if err != nil {
		log.Fatal("Failed to parse RFC 2849 from input file: "+os.Args[1], err)
	}

	lineCount := 0
	for i, e := range l.AllEntries() {

		mail := strings.Join(e.GetAttributeValues("mail"), ",")

		for j, certBytes := range e.GetRawAttributeValues("userCertificate;binary") {
			var csvFields [8]string
			csvFields[0] = fmt.Sprintf("(%d/%d)", i+1, j+1)
			csvFields[1] = mail
			csvFields[2] = e.DN

			cert, err := x509.ParseCertificate(certBytes)
			if err != nil {
				// When the err is "asn1: structure error: integer not minimally-encoded"
				// See: https://github.com/google/certificate-transparency-go/issues/340
				csvFields[3] = fmt.Sprintf("Can't parse cert: %s\n", err)
			} else {
				csvFields[3] = cert.Issuer.String()
				csvFields[4] = cert.Subject.String()
				csvFields[5] = cert.SerialNumber.String()
				csvFields[6] = cert.NotBefore.String()
				csvFields[7] = cert.NotAfter.String()
			}
			csvWriter.Write(csvFields[:])
			lineCount = lineCount + 1
		}
	}
	fmt.Printf("%d records written to CSV file\n", lineCount)
}
