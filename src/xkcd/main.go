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

// xkcd

// We start with the current day, then follow the Previous link until we find
// an image we've already saved before. We also save both the 1x and 2x
// versions of images, when available.

var start = "https://xkcd.com/"
var file = regexp.MustCompile("//imgs.xkcd.com/comics/([^\"]+\\.png)")
var link = regexp.MustCompile("rel=\"prev\" href=\"/([0-9]+/)\"")

func main() {
	os.MkdirAll("comics/xkcd", os.ModePerm)

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

		// Find comic images
		files := file.FindAllStringSubmatch(s, 2)
		for i := range files {
			path := "comics/xkcd/" + files[i][1]
			imgurl := "https://imgs.xkcd.com/comics/" + files[i][1]

			// Load image and write to file
			if _, err := os.Stat(path); os.IsNotExist(err) {
				// Image does not exist locally
				fmt.Println("Downloading:", files[i][1])
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
		}

		// Find link to previous comic
		links := link.FindStringSubmatch(s)
		url = start + links[1]

		// Wait a bit
		time.Sleep(500 * time.Millisecond)
	}
}
