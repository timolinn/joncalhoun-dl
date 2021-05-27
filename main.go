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
	"path/filepath"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/publicsuffix"
)

var email = flag.String("email", "", "your account email")
var password = flag.String("password", "", "your account password")
var course = flag.String("course", "gophercises", "course name")
var outputdir = flag.String("output", "", "output directory\n ie. where the course videos are saved")
var cachelocation = flag.String("cache", "", "specifies where the cache will be saved\nDefaults to a joncalhoun-dl-cache folder within the output directoy")
var help = flag.Bool("help", false, "prints this output")

// this will be used by youtube-dl binary to download video
var referer = "https://courses.calhoun.io"

var courses = map[string]string{
	"testwithgo":       "https://courses.calhoun.io/courses/cor_test",
	"gophercises":      "https://courses.calhoun.io/courses/cor_gophercises",
	"algorithmswithgo": "https://courses.calhoun.io/courses/cor_algo",
	"webdevwithgo":     "https://courses.calhoun.io/courses/cor_webdev",
}
var delayDuration = 5

// ClientOption is the type of constructor options for NewClient(...).
type ClientOption func(*http.Client) error

func checkError(err error) {
	if err != nil {
		log.Fatalf("[joncalhoun-dl]:[Error]: %v.\nif you think this is unusual please create an issue %s\n",
			err,
			"https://github.com/timolinn/joncalhoun-dl/issues/new")
	}
}

// NewClient constructs anew client which can make requests
// to course website
func NewClient(options ...ClientOption) (*http.Client, error) {
	// Cookiejar provides automatic cookie management
	// that would normally be accessed only via the browser
	opts := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&opts)
	checkError(err)
	c := &http.Client{Jar: jar}
	for _, option := range options {
		err := option(c)
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

// WithTransport configures the client to use a different transport
func WithTransport(fn RoundTripperFunc) ClientOption {
	return func(client *http.Client) error {
		client.Transport = RoundTripperFunc(fn)
		return nil
	}
}

func main() {
	// Parse commandline options
	flag.Parse()

	// Print usage options
	if *help {
		flag.PrintDefaults()
		return
	}

	// validate input
	validateInput()

	client, err := NewClient()
	checkError(err)

	// Login
	signin(client)

	location := *outputdir + "/%(title)s.%(ext)s"
	if *outputdir == "" {
		cwd, err := os.Getwd()
		checkError(err)
		*outputdir = cwd + "/" + *course
		location = *outputdir + "/%(title)s.%(ext)s"
	}
	fmt.Printf("[joncalhoun-dl]: output directory is %s\n", *outputdir)

	// do some chores
	setup()

	// fetch video urls
	videoURLs := getURLs(client)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for i, videoURL := range videoURLs {
		if videoURL != "" {
			fmt.Printf("[joncalhoun-dl]: downloading lesson 0%d of %s\n", i+1, *course)
			fmt.Printf("[joncalhoun-dl]:[exec]: youtube-dl %s --referer %s -o %s\n", videoURL, referer, location)
			cmd := exec.CommandContext(ctx, "youtube-dl", videoURL, "--referer", referer, "-o", location)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Start(); err != nil {
				log.Fatal(err)
			}
			if err := cmd.Wait(); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("[joncalhoun-dl]: downloaded lesson 0%d\n", i+1)
		} else {
			fmt.Printf("[joncalhoun-dl]: Page for lesson 0%d does not have an embedded video \n", i+1)
		}
	}

	os.Chdir(*outputdir)

	files, _ := filepath.Glob("*.mp4")

	for _, file := range files {

		slc := strings.Split(file, "-")
		for i := range slc {
			slc[i] = strings.TrimSpace(slc[i])
		}

		lesson := strings.Join(slc[:2], " ")
		video := strings.Join(slc[2:], " ")

		// Create lesson directory if it does not exist yet
		if !dirExists(*outputdir + "/" + lesson) {
			err := os.Mkdir(*outputdir+"/"+lesson, 0755)
			checkError(err)
		}

		fmt.Printf("[joncalhoun-dl]: Moving %s to %s/%s\n", file, lesson, video)

		os.Rename(file, lesson+"/"+video)
	}

	fmt.Println("Done! 🚀")
}

func validateInput() {
	if !isSupported(*course) {
		err := errors.New("course not supported yet")
		checkError(err)
	}

	if *email == "" || *password == "" {
		checkError(errors.New("try adding: --email=jon@examp.com --password=12345' to the command"))
	}
}

func setup() {
	// create output directory if it does not exist yet
	if !dirExists(*outputdir) {
		err := os.Mkdir(*outputdir, 0755)
		checkError(err)
	}

	if *cachelocation == "" {
		*cachelocation = *outputdir + "/" + "joncalhoun-dl-cache"
	}

	// create cache location if it does not exist
	if !dirExists(*cachelocation) {
		err := os.Mkdir(*cachelocation, 0755)
		checkError(err)
	}
}

func signin(client *http.Client) {
	// Login and create session

	fmt.Println("[joncalhoun-dl]: signing in...")
	_, err := client.PostForm("https://courses.calhoun.io/signin", url.Values{
		"email":    {*email},
		"password": {*password},
	})
	checkError(err)
	fmt.Println("[joncalhoun-dl]: sign in successful")
}

func getCourseHTML(url string, client *http.Client) {
	// Make a Get Request to the course URL and fetch the HTML
	// user must be logged in
	fmt.Printf("[joncalhoun-dl]: fetching data for %s...\n", url)
	res, err := client.Get(url)
	checkError(err)
	defer res.Body.Close()

	// Write data to file
	saveHTMLContent(*course+".html", res.Body)
}

func getURLs(client *http.Client) []string {
	fmt.Printf("[joncalhoun-dl]: fetching video urls for %s\n", *course)
	var urls []string
	var file *os.File
	var err error

	// check if course page is cached
	if isCached(*course + ".html") {
		fmt.Printf("[joncalhoun-dl]: loading %s data from cache \n", *course)
		file, err = loadFromCache(*course + ".html")
		checkError(err)
	} else {
		// fecth from remote if not cached
		fmt.Printf("[joncalhoun-dl]: fetching %s data from remote\n", *course)
		res, err := client.Get(courses[*course])
		checkError(err)
		defer res.Body.Close()

		// cache raw HTML data
		getCourseHTML(courses[*course], client)
		file, err = loadFromCache(*course + ".html")
		checkError(err)
	}

	doc, err := goquery.NewDocumentFromReader(file)
	checkError(err)

	// parses the HTML tree to extract url
	// where the lesson video is located
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		switch *course {
		case "testwithgo":
			// each lesson link should contain this substring
			// else ignore
			if strings.Contains(href, "/lessons/les_twg") {
				urls = append(urls, "https://courses.calhoun.io"+href)
			}
		case "gophercises":
			// each lesson link should contain this substring
			// else ignore
			if strings.Contains(href, "/lessons/les_goph") {
				urls = append(urls, "https://courses.calhoun.io"+href)
			}
		case "webdevwithgo":
			if strings.Contains(href, "/lessons/les_wd") {
				urls = append(urls, "https://courses.calhoun.io"+href)
			}
		case "advancedwebdevwithgo":
			log.Fatal("'Advanced Web Development with Go' not supported yet")
		case "algorithmswithgo":
			if strings.Contains(href, "/lessons/les_algo") {
				urls = append(urls, "https://courses.calhoun.io"+href)
			}
		default:
			log.Fatal("course not supported yet. feel free to send a pull request")
		}
	})

	videoURLs := []string{}
	for _, url := range urls {
		videoURLs = append(videoURLs, getVideoURL(url, client))
		// we don't want to send too many requests in a short time
		// this naively simulates human behaviour
		fmt.Printf("[joncalhoun-dl]: waiting 5 seconds\n")
		time.Sleep(time.Duration(delayDuration) * time.Second)
	}
	return videoURLs
}

func getVideoURL(url string, client *http.Client) string {
	fmt.Printf("[joncalhoun-dl]: fetching video url for lesson %s\n", url)
	var videoID string
	var file *os.File
	var err error

	// check cache for existing webpage
	name := strings.Split(url, "/")[4]
	filename := name + ".html"
	if isCached(filename) {
		fmt.Printf("[joncalhoun-dl]: loading %s from cache\n", name)
		file, err = loadFromCache(filename)
		checkError(err)

		// no need to delay when loading from cash
		delayDuration = 0
	} else {
		// fetch web page where video lives
		fmt.Printf("[joncalhoun-dl]: fetching %s from remote\n", filename)
		res, err := client.Get(url)
		checkError(err)
		defer res.Body.Close()

		// To provide caching support we save the resulting
		// html in the cache folder
		saveHTMLContent(filename, res.Body)
		file, err = loadFromCache(filename)
		delayDuration = 5
	}

	// convert return data to parsable HTML Document
	doc, err := goquery.NewDocumentFromReader(file)
	checkError(err)
	iframe := doc.Find("iframe")
	videoID, _ = iframe.Attr("src")
	fmt.Printf("[joncalhoun-dl]:[video ID] %s\n", videoID)
	return videoID
}

func saveHTMLContent(filename string, r io.Reader) {
	f, err := os.Create(*cachelocation + "/" + filename)
	checkError(err)
	defer f.Close()
	filewriter := bufio.NewWriter(f)
	_, err = filewriter.ReadFrom(r)
	checkError(err)

	filewriter.Flush()
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func isCached(name string) bool {
	if fileExists(*cachelocation + "/" + name) {
		return true
	}
	return false
}

func loadFromCache(name string) (*os.File, error) {
	return os.OpenFile(*cachelocation+"/"+name, os.O_RDWR, 0666)
}

func isSupported(coursename string) bool {
	if courses[coursename] != "" {
		return true
	}
	return false
}
