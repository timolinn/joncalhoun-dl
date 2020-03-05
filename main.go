package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/publicsuffix"
)

var email = flag.String("email", "", "your email")
var password = flag.String("password", "", "your password")
var course = flag.String("course", "gophercises", "course name")

// this will be used by youtube-dl binary to download video
var referer = "https://courses.calhoun.io"

var courses = map[string]string{
	"testwithgo":  "https://courses.calhoun.io/courses/cor_test",
	"gophercises": "https://courses.calhoun.io/courses/cor_gophercises",
	"algorithms":  "https://courses.calhoun.io/courses/cor_algo",
	"webdevelopmentwithgo": "https://courses.calhoun.io/courses/cor_webdev",
}
var delayDuration = 5

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Parse commandline options
	flag.Parse()

	// Cookiejar provides automatic cookie management
	// that would normally be accessed only via the browser
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	checkError(err)

	client := &http.Client{Jar: jar}

	// Login
	signin(client)

	// Visit selected course page and fetch video urls
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	videoURLs := getURLs(client)
	for i, videoURL := range videoURLs {
		fmt.Printf("[joncalhoun.io]: downloading lesson 0%d via %s\n", i+1, videoURL)
		fmt.Printf("[exec]: youtube-dl %s --referer %s\n", videoURL, referer)
		cmd := exec.CommandContext(ctx, "youtube-dl", videoURL, "--referer", referer, "-o",
			"./videos/"+*course+"/%(title)s.%(ext)s")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}
		if err := cmd.Wait(); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("[joncalhoun.io]: downloaded lesson 0%d\n", i+1)
	}
	fmt.Println("Done! 🚀")
}

func signin(client *http.Client) {
	// Login and create session
	if *email == "" || *password == "" {
		log.Fatal(errors.New("[Error] try: 'go run main.go --email=jon@examp.com --password=12345'"))
	}

	fmt.Println("[joncalhoun.io]: signing in...")
	_, err := client.PostForm("https://courses.calhoun.io/signin", url.Values{
		"email":    {*email},
		"password": {*password},
	})
	checkError(err)
	fmt.Println("[joncalhoun.io]: sign in successful")
}

func getCourseHTML(url string, client *http.Client) {
	// Make a Get Request to the course URL
	// user must be logged in
	fmt.Printf("[joncalhoun.io]: fetching data for %s...\n", url)
	res, err := client.Get(url)
	checkError(err)
	defer res.Body.Close()

	// Write raw data to file
	saveHTMLContent(*course+".html", res.Body)
}

func getURLs(client *http.Client) []string {
	fmt.Printf("[joncalhoun.io]: fetching video urls for %s\n", *course)
	var urls []string
	var file *os.File
	var err error

	// check if course page is cached
	if fileExists(*course + ".html") {
		fmt.Printf("[joncalhoun.io]: loading course page from cache: %s.html\n", *course)
		file, err = os.OpenFile(*course+".html", os.O_RDWR, 0666)
		checkError(err)
	} else {
		// fecth from remote if not cached
		res, err := client.Get(courses[*course])
		checkError(err)
		defer res.Body.Close()

		// cache raw HTML data
		getCourseHTML(courses[*course], client)
		file, err = os.OpenFile(*course+".html", os.O_RDWR, 0666)
		checkError(err)
	}

	doc, err := goquery.NewDocumentFromReader(file)
	checkError(err)

	// parses the HTML tree to extract videos pages
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		switch *course {
		case "testwithgo":
			if strings.Contains(href, "/lessons/les_twg") {
				urls = append(urls, "https://courses.calhoun.io"+href)
			}
		case "gophercises":
			if strings.Contains(href, "/lessons/les_goph") {
				urls = append(urls, "https://courses.calhoun.io"+href)
			}
		case "algorithms":
			log.Fatal("Algorithms not supported yet")
		default:
			log.Fatal("course not supported yet. feel free to send a pull request")
		}
	})

	videoURLs := []string{}
	for _, url := range urls {
		videoURLs = append(videoURLs, getVideoURL(url, client))
		fmt.Printf("[joncalhoun.io]: waiting 5 seconds\n")
		time.Sleep(time.Duration(delayDuration) * time.Second)
	}
	fmt.Println("Nothing here")
	return videoURLs
}

func getVideoURL(url string, client *http.Client) string {
	fmt.Printf("[joncalhoun.io]: fetching video url for lesson %s\n", url)
	var videoID string
	var file *os.File
	var err error

	// check cache for existing webpage
	name := strings.Split(url, "/")[4]
	filename := "webpage/" + name + ".html"
	if fileExists(filename) {
		fmt.Printf("[joncalhoun.io]: loading cache from %s\n", filename)
		file, err = os.OpenFile(filename, os.O_RDWR, 0666)
		checkError(err)

		// no need to delay when loading from cash
		delayDuration = 0
	} else {
		// fetch web page where video lives
		res, err := client.Get(url)
		checkError(err)
		defer res.Body.Close()

		// To provide caching support we save the resulting
		// html in the webpage folder
		saveHTMLContent(filename, res.Body)
		file, err = os.OpenFile(filename, os.O_RDWR, 0666)
		checkError(err)
		delayDuration = 5
	}

	// convert return data to parsable HTML Document
	doc, err := goquery.NewDocumentFromReader(file)
	checkError(err)
	iframe := doc.Find("iframe")
	videoID, _ = iframe.Attr("src")
	fmt.Printf("[joncalhoun.io]:[video ID] %s\n", videoID)
	return videoID
}

func saveHTMLContent(filename string, r io.Reader) {
	f, err := os.Create(filename)
	checkError(err)
	defer f.Close()
	filewriter := bufio.NewWriter(f)
	_, err = filewriter.ReadFrom(r)
	checkError(err)
	fmt.Printf("[joncalhoun.io]: web page data written to %s\n", filename)

	filewriter.Flush()
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
