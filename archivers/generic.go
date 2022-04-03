package archivers

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"
)

// Generic archiver supporting ComicPress, etc.
func Generic(startURL string, dir string, fileMatch *regexp.Regexp, filePrefix string, prevLinkMatch *regexp.Regexp, skipExisting bool) {
	os.MkdirAll("comics/"+dir, os.ModePerm)

	fmt.Println(startURL)

	url := startURL
	for {
		// Load page HTML
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(dir, "Failed to load page:", err)
			os.Exit(1)
		}
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		s := buf.String()

		// Find comic image
		files := fileMatch.FindStringSubmatch(s)
		if len(files) < 1 {
			fmt.Println(dir, "No comic image found")
			return
		}
		imgUrl := filePrefix + files[1]
		basename := basenameMatch.FindStringSubmatch(files[1])
		path := "comics/" + dir + "/" + basename[1]

		// Download image
		dlErr := downloadFileWait(basename[1], path, imgUrl, 500*time.Millisecond)
		if dlErr != nil {
			if dlErr.Error() == "file exists" {
				if !skipExisting {
					fmt.Println(dir, "File exists:", path)
					return
				}
			} else {
				fmt.Println(dir, "Error:", dlErr.Error())
				return
			}
		}

		// Find link to previous comic
		links := prevLinkMatch.FindStringSubmatch(s)
		if len(links) < 1 {
			fmt.Println(dir, "No previous URL found")
			return
		}
		if protocolMatch.MatchString(links[1]) {
			url = links[1]
		} else {
			url = startURL + links[1]
		}

		// Wait a bit
		time.Sleep(500 * time.Millisecond)
	}
}

func GenericCustomStart(startURL string, startMatch *regexp.Regexp, dir string, fileMatch *regexp.Regexp, filePrefix string, prevLinkMatch *regexp.Regexp, skipExisting bool) {
	// Load start page HTML
	resp, err := http.Get(startURL)
	if err != nil {
		fmt.Println(dir, "Failed to load page:", err)
		os.Exit(1)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	s := buf.String()

	// Find comic start page
	start := startMatch.FindStringSubmatch(s)
	if len(start) < 1 {
		fmt.Println(dir, "No start URL found")
		os.Exit(1)
	}

	// Start comic download
	Generic(start[1], dir, fileMatch, filePrefix, prevLinkMatch, skipExisting)
}
