package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"log"
	"regexp"
)

func main() {
	hostname1 := myhostname1()
	fmt.Println(hostname1)

	hostname2 := myhostname2()
	fmt.Println(hostname2)
}

func myhostname1() string {
	res, err := http.Get("https://www.displaymyhostname.com/")
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	var hostnamePattern = regexp.MustCompile(`<div class="hostname">(.*)<\/div>`)
	matches := hostnamePattern.FindStringSubmatch( string(body))
	fmt.Println( string(body))
	fmt.Println(matches[1])
	return "cica.hu"
}

func myhostname2() string {
	res, err := http.Get("https://myhostname.net/")
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	var hostnamePattern = regexp.MustCompile(` <span id="curHostname" class="notranslate">(.*)<\/span>`)
	matches := hostnamePattern.FindStringSubmatch( string(body))
	fmt.Println( string(body))
	fmt.Println(matches[1])
	return "cica.hu"
}
