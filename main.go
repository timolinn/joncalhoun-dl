package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func main() {
	res, err := http.PostForm("https://courses.calhoun.io/signin", url.Values{
		"email":    {"timothyonyiuke@gmail.com"},
		"password": {"cyberelf**20"},
	})
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
}

func temp() {
	client := &http.Client{}
	// header := http.Header{}
	// data, err := ioutil.ReadFile("cookies.txt")
	// cooks, err := cookiestxt.Parse(strings.NewReader(string(data)))
	// header.Add("Cookie")
	req := http.Request{
		// Header: header,
		Method: "POST",
		URL: &url.URL{
			Host:   "courses.calhoun.io",
			Path:   "/signin",
			Scheme: "https",
		},
		PostForm: url.Values{
			"email":    []string{"timothyonyiuke@gmail.com"},
			"password": []string{"cyberelf**20"},
		},
	}
	// req.AddCookie(cooks[0])
	fmt.Println(req.Cookies())
	// res, err := http.Get("https://courses.calhoun.io/lessons/les_goph_128")
	// res, err := http.NewRequestWithContext(req.Context, "GET", "https://courses.calhoun.io/lessons/les_goph_128", )
	// err := req.Write(os.Stdout)
	res, err := client.Do(&req)
	// req.AddCookie(getCookie())
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
}
