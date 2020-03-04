package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var email = flag.String("email", "", "your email")
var password = flag.String("password", "", "your password")
var course = flag.String("course", "testwithgo", "course name")

// this will be used by youtube-dl binary to download video
var referer = "https://courses.calhoun.io"

var courses = map[string]string{
	"testwithgo":  "https://courses.calhoun.io/courses/cor_test",
	"gophercises": "https://courses.calhoun.io/courses/cor_gophercises",
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Parse commanline options
	flag.Parse()

	// Cookiejar provides automatic cookie management
	// that would normally be aaccessed only via the browser
	// options := cookiejar.Options{
	// 	PublicSuffixList: publicsuffix.List,
	// }
	// jar, err := cookiejar.New(&options)
	// checkError(err)

	// client := &http.Client{Jar: jar}

	// Login
	// res, client := signin(client)

	// Visit selected course page and fetch video paths
	// There are currently 4 different path pattern
	// `/les_twg_les_01` => for tESTwITHgO tests videos
	// `/les_form_les_01` => for tESTwITHgO form project videos
	// `/les_stripe_les_01` => for tESTwITHgO stripe project videos
	// `/les_swag_les_01` => for tESTwITHgO swag project videos
	// res, err := client.Get(courses[*course])
	// checkError(err)
	getURLs()
}

func signin(client *http.Client) (*http.Response, *http.Client) {
	// Login and create session
	if *email == "" || *password == "" {
		log.Fatal(errors.New("[Error] try: 'go run main.go --email=jon@examp.com --password=12345'"))
	}

	res, err := client.PostForm("https://courses.calhoun.io/signin", url.Values{
		"email":    {*email},
		"password": {*password},
	})
	checkError(err)

	body, err := ioutil.ReadAll(res.Body)
	checkError(err)
	res.Body.Close()
	fmt.Println(string(body))
	return res, client
}

func getCourseHTML(name string, res io.Reader) {
	f, err := os.Create(name)
	checkError(err)
	defer f.Close()
	writer := bufio.NewWriter(f)
	num, err := writer.ReadFrom(res)
	checkError(err)
	fmt.Printf("Wrote %d bytes", num)

	writer.Flush()
}

func getURLs() []string {
	var urls []string
	f, err := os.OpenFile("testwithgo.html", os.O_RDWR, 066)
	doc, err := goquery.NewDocumentFromReader(f)
	checkError(err)
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		if strings.Contains(href, "/lessons/les_twg") {
			urls = append(urls, "https://courses.calhoun.io"+href)
		}
	})
	return urls
}
