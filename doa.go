package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"time"
)

// Dumbing of Age

// We start with the current day, then follow the Previous link until we find
// an image we've already saved before.

var start = "http://www.dumbingofage.com/"
var file = regexp.MustCompile("/comics/(.+\\.png)")
var link = regexp.MustCompile("href=\"(http://www.dumbingofage.com/[0-9a-zA-Z/-]+)\" class=\"navi navi-prev\"")

func main() {
	os.MkdirAll("comics/doa", os.ModePerm)

	url := start
	for {
		// Load page HTML
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Failed to load page:", err)
			os.Exit(1)
		}
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		s := buf.String()

		// Find comic image
		files := file.FindStringSubmatch(s)
		path := "comics/doa/" + files[1]
		imgurl := "http://www.dumbingofage.com/comics/" + files[1]

		// Load image and write to file
		if _, err := os.Stat(path); os.IsNotExist(err) {
			// Image does not exist locally
			fmt.Println("Downloading:", files[1])
			imgresp, err := http.Get(imgurl)
			if err != nil {
				fmt.Println("Failed to load image:", err)
				os.Exit(1)
			}

			buf := new(bytes.Buffer)
			buf.ReadFrom(imgresp.Body)
			err = ioutil.WriteFile(path, buf.Bytes(), 0644)
			if err != nil {
				fmt.Println("Failed to write image:", err)
				os.Exit(1)
			}
		} else {
			fmt.Println("File exists:", path)
			os.Exit(0)
		}

		// Find link to previous comic
		links := link.FindStringSubmatch(s)
		url = links[1]

		// Wait a bit
		time.Sleep(500 * time.Millisecond)
	}
}
