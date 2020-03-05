package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveHTMLContent(t *testing.T) {
	t.Run("saveHTMLContent: should saves html file to filesystem", func(t *testing.T) {
		filename := "temp.html"
		r := strings.NewReader(`<!DOCTYPE html>
		<html lang="en">
		  <head>
			<meta charset="utf-8" />

			<link rel="shortcut icon" href="/favicon.ico" />
			<meta name="viewport" content="width=device-width, initial-scale=1" />
			<meta name="theme-color" content="#000000" />

			<link
			  href="https://fonts.googleapis.com/css?family=Lato|Satisfy|Kristi&display=swap"
			  rel="stylesheet"
			/>
			<link href="/assets/styles.css" rel="stylesheet">

			<title>Gophercises | courses.calhoun.io</title>
		  </head>
		  <body class="bg-grey-100">
		  <div></div>
		  </body>
		  </html>
		  `)
		saveHTMLContent(filename, r)

		_, err := os.Stat(filename)
		if err != nil {
			t.Error()
			return
		}
		assert.False(t, os.IsNotExist(err))
		if err := os.Remove(filename); err != nil {
			t.Error()
			return
		}
	})
}

func TestGetCourseHTML(t *testing.T) {
	handler := func(r *http.Request) (*http.Response, error) {
		body := `<!DOCTYPE html>
		<html lang="en">
		  <head>
			<meta charset="utf-8" />

			<link rel="shortcut icon" href="/favicon.ico" />
			<meta name="viewport" content="width=device-width, initial-scale=1" />
			<meta name="theme-color" content="#000000" />

			<link
			  href="https://fonts.googleapis.com/css?family=Lato|Satisfy|Kristi&display=swap"
			  rel="stylesheet"
			/>
			<link href="/assets/styles.css" rel="stylesheet">

			<title>Gophercises | courses.calhoun.io</title>
		  </head>
		  <body class="bg-grey-100">


		  <div class="w-full mb-4 pt-8">
			<div>
			  <h3 class="text-grey-600 border-b border-grey-200 py-1 text-2xl flex items-baseline mx-6">
				<div>Quiz Game</div>
			  </h3>

				<div class="text-grey-600 px-6 py-2 markdown">
				  <p>Create a program to run timed quizes via the command line.</p>

				</div>

			  <div class="flex flex-wrap px-2 pt-4 pb-12">

				  <a href="/lessons/les_goph_01">
					<div class="w-64 mb-4 mx-4">
					  <div class="w-full h-144px">
						<img alt="Thumbnail for Overview" src="/assets/img/thumbs/les_goph_01.png" class="rounded max-h-full" />
					  </div>
					  <div class="px-2 py-2">
						<span class="text-grey-700 font-sans no-underline">
						  Overview
						</span>
					  </div>
					</div>
				  </a>

				  <a href="/lessons/les_goph_02">
					<div class="w-64 mb-4 mx-4">
					  <div class="w-full h-144px">
						<img alt="Thumbnail for Solution - Part 1" src="/assets/img/thumbs/les_goph_02.png" class="rounded max-h-full" />
					  </div>
					  <div class="px-2 py-2">
						<span class="text-grey-700 font-sans no-underline">
						  Solution - Part 1
						</span>
					  </div>
					</div>
				  </a>

				  <a href="/lessons/les_goph_03">
					<div class="w-64 mb-4 mx-4">
					  <div class="w-full h-144px">
						<img alt="Thumbnail for Solution - Part 2" src="/assets/img/thumbs/les_goph_03.png" class="rounded max-h-full" />
					  </div>
					  <div class="px-2 py-2">
						<span class="text-grey-700 font-sans no-underline">
						  Solution - Part 2
						</span>
					  </div>
					</div>
				  </a>

			</div>
			</div>
			</div>`

		return &http.Response{
			Body: ioutil.NopCloser(strings.NewReader(body)),
		}, nil
	}
	t.Run("getCourseHTML: fetches the course HTML main page, this page contains links to individual lessons that make up the course", func(t *testing.T) {
		client, _ := NewClient(WithTransport(handler))
		// set course name
		*course = *course + "_test"
		getCourseHTML(courses["gophercises"], client)
		assert.True(t, fileExists(*course+".html"))
		os.Remove(*course + ".html")
	})
}

func TestFileExists(t *testing.T) {
	t.Run("fileExists: should check if a file exists when given a path", func(t *testing.T) {
		assert.True(t, fileExists("README.md"))
		assert.False(t, fileExists("nonExistingFile"))
	})
}

func TestDirExists(t *testing.T) {
	t.Run("dirExists: should check if a directory exists when given a path", func(t *testing.T) {
		assert.True(t, dirExists("webpage"))
		assert.False(t, dirExists("nonExistingDir"))
	})
}

func init() {
	os.Setenv("APP_ENV", "test")
}
