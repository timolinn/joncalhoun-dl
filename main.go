package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"golang.org/x/net/publicsuffix"
)

var email = flag.String("email", "", "email")
var password = flag.String("password", "", "password")

func main() {
	flag.Parse()

	if *email == "" || *password == "" {
		log.Fatal(errors.New("[Error] try: 'go run main.go -email=jon@examp.com -password=12345'"))
	}
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{Jar: jar}
	res, err := client.PostForm("https://courses.calhoun.io/signin", url.Values{
		"email":    {*email},
		"password": {*password},
	})
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	// fmt.Println(res.Cookies())
	fmt.Println(string(body))
}
