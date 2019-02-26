// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/tls"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bukodi/go-playgroud/swaggerui"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

var templates *template.Template

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func swaggeruiHandler(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/swagger-ui/") {
		http.NotFound(w, r)
		return
	}

	assetName := r.URL.Path[len("/swagger-ui/"):]
	bytes, err := swaggerui.Asset(assetName)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if strings.HasSuffix(assetName, ".css") {
		w.Header().Set("Content-Type", "text/css")
	}

	w.Write(bytes)
	//r.URL.Path
}

var certPem = []byte(`-----BEGIN CERTIFICATE-----
MIIC+jCCAeKgAwIBAgIRAJ+6uxaxKahykPSLdsLZ3hAwDQYJKoZIhvcNAQELBQAw
EjEQMA4GA1UEChMHQWNtZSBDbzAeFw0xODAzMTkxNDU2NDNaFw0xOTAzMTkxNDU2
NDNaMBIxEDAOBgNVBAoTB0FjbWUgQ28wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAw
ggEKAoIBAQCyuwtH+jxrtihDhVDz7XurFI05DycjXeHgyjYvUzE2o6SK09W0yrqe
lAySGPAXYGKldQtDiFHuOFrM8qAKVcrb7I/+O08gJ9Y5WR5jLer8UlDhpuRQBosW
Ts4ReeeXbHyu5sV2ODSW9hGs6A8OIyALXrgZ05aa9oI4NHd0iRv8WsfbCRGD3zNV
I6RkGyvlWnuy3FVx/xJdir3l26+kYNeb6rMQ79+gN7Vr7jCSM3k4wsLaCocnbLQh
zl1jgrJBUzKfa7HxUOPCXucQvuEi/4SO4retyHnuDdqfMTGznjRdlzdb+CnOaheC
kD4NThFNOCeatezjbsMio6q7B4oMhVyjAgMBAAGjSzBJMA4GA1UdDwEB/wQEAwIF
oDATBgNVHSUEDDAKBggrBgEFBQcDATAMBgNVHRMBAf8EAjAAMBQGA1UdEQQNMAuC
CWxvY2FsaG9zdDANBgkqhkiG9w0BAQsFAAOCAQEAAvQgN/iIK70rwuEuhEMZPe8Q
0LF3cPVtG2AyPLaQ+AaKkFbRhmbD45YzYYCb5E4XcAHnUcqyYWT8EzLVx5zNB6FW
j1YgtCavDygP9Ecdc/pQ1ocHPUcCbs9H7slYf2tYICJ1zcOYYxbndri2SZcj1CiW
JXoKDEXY3zb582J8rL+01pl4n5/HDKz/Xe2rizZR9HzGnIJzOuQiAP4R8cKZFr1J
LAQ0GKxSF1p97LKLgvDt4B9g34EEmqGvTg8WrhNMXUso3qYeubqKtQ8+VDru9aPn
XC5UOwcLNFvbUuOtHyZoiB3QWDZEAxuxFvYECBuUz7xi6vbEhfvYYxWQCXvt1w==
-----END CERTIFICATE-----`)

var keyPem = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAsrsLR/o8a7YoQ4VQ8+17qxSNOQ8nI13h4Mo2L1MxNqOkitPV
tMq6npQMkhjwF2BipXULQ4hR7jhazPKgClXK2+yP/jtPICfWOVkeYy3q/FJQ4abk
UAaLFk7OEXnnl2x8rubFdjg0lvYRrOgPDiMgC164GdOWmvaCODR3dIkb/FrH2wkR
g98zVSOkZBsr5Vp7stxVcf8SXYq95duvpGDXm+qzEO/foDe1a+4wkjN5OMLC2gqH
J2y0Ic5dY4KyQVMyn2ux8VDjwl7nEL7hIv+EjuK3rch57g3anzExs540XZc3W/gp
zmoXgpA+DU4RTTgnmrXs427DIqOquweKDIVcowIDAQABAoIBAEA436PcceuOR8eD
VwRfeEmQF/LB1CFsMabxYij9LrjgclaEKc1N72Ld9eplVZhAxRGJDiQVDsOXsmDk
acds7ni59z+2FgeK6PowYK/opwBFn9SFgJKU11OLu5YiBKqvi9nZCGkjZxo7jzxf
IlHFI1WSknqNQheqrj399FKGlezGJCFTpCw3DzKTudbExCCVMO1d//8Y2rKoRrz5
iqSsekSqEeCnTgv22w4BJsM6sOVHZSJ02bYCm5ps4IGYy9KSOaJc979dSJKqAbTW
1vr+FHUrhKKi/sB/LjdmaGc5u//3JumssVvh7wXR8Ph74kDgsDrIE4bzlL1Uaei3
J7ta30kCgYEAzUEt13I20g84MH5ZzSSf1h0ptg2FF2V9vXa3886gxwjpka9WDRbV
PCaJsjvRVFrqjnSCFJFJwQyg7zveFitJzddhufMNjqnY/V84Dbnx4PgkUeOsumTR
zw/7oq1EEO7dQqjivgDddPmRNf0xVWB5lBlYsnqSSUs6pSlkUFM3H20CgYEA3ush
S6jkUWI9SOL6B9jAG/72Qfo0Jci4DUJWXShTHeD00Y8ScyFZyaFyDVc0Ywbpbj4X
5hpKi7Xja+6Zj22lKUoG95Zed4vhGvDHXl1bAK9YT4YrQb4y1GefDXJzipcZdV63
8QJhBYzo/5Q3tS/4lRkm5EzuvhRRibQyFi5hEk8CgYEAsLXn6K/tYKY3wxBU8hgR
AD81VQaIyh7XxZF7SA1iQFl89a9Vz4kT5mhrbiavzwdDH4hRIbIAJJNhzvXk+4Mj
VHOVMIl/5451QZaD5NVs2Dnq0xHH+OWp+LIS+/hePJHZrnVGWTzXbMkcarXkjlOz
+Hxl76s1XKLHB8D+G2W5dHUCgYACoDpwLbkizEl0hlfzp7X7nnFALbZXi5m/bjye
NE9mVrQLk+ffu1DXczNovNI9KGOvjMOzTjP6mVXoe5MLgXsklV6no/nQ5rDsJFH0
5pyf0XD03tu7loX6wo25FtQNmeIO4+K+0AxciGBmQlS1qa7/8p/mqJFXY93iBWFh
qYIzOQKBgQDBz1DFJtiZxgUgPVfeMaE7GqpqorD8CVSQ4CLmTjYPsCbrfxtRiZ60
SCMwJ5zXLhcio/E+XqEIk2hrseVVOeILpC3+P2/klGFJjxML7Tth1LhyztAzCWiB
b1ARF33V/sYZ9TAcR2RfXEEVps4TAU2yYsY3a6EhdBr7Osrumk0aBA==
-----END RSA PRIVATE KEY-----`)

func main() {

	basePath, _ := filepath.Abs(".")
	fmt.Println(basePath)

	templates = template.Must(template.ParseFiles("edit.html", "view.html"))

	fmt.Println(swaggerui.Asset("index.html"))

	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/swagger-ui/", swaggeruiHandler)

	fmt.Println("Started on port ", 8080)

	server := &http.Server{Addr: ":10443", Handler: nil}

	cer, err := tls.X509KeyPair(certPem, keyPem)
	if err != nil {
		fmt.Println(err)
		return
	}
	tlsCfg := &tls.Config{Certificates: []tls.Certificate{cer}}

	conn, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Fatal(err)
	}

	tlsListener := tls.NewListener(conn, tlsCfg)
	server.Serve(tlsListener)

}
