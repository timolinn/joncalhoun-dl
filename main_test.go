package main

import (
	"os"
	"strings"
	"testing"
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
		if os.IsNotExist(err) {
			t.Error("failed to save file")
			return
		}
		if err := os.Remove(filename); err != nil {
			t.Error()
			return
		}
	})
}

func TestFileExists(t *testing.T) {
	t.Run("fileExists: should check if a file exists when given a path", func(t *testing.T) {
		if !fileExists("README.md") {
			t.Error()
			return
		}

		if fileExists("non-existing-file.md") {
			t.Error()
			return
		}
	})
}
