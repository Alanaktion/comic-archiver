package archivers

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"
)

// Generic archiver supporting ComicPress, etc.
func Generic(startURL string, dir string, fileMatch *regexp.Regexp, filePrefix string, prevLinkMatch *regexp.Regexp, skipExisting bool) error {
	os.MkdirAll("comics/"+dir, os.ModePerm)

	fmt.Println(startURL)

	url := startURL
	for {
		// Load page HTML
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(dir, "Failed to load page:", err)
			return err
		}
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		s := buf.String()

		// Find comic image
		files := fileMatch.FindStringSubmatch(s)
		if len(files) < 1 {
			fmt.Println(dir, "No comic image found")
			return errors.New("no comic image found")
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
					return nil
				}
			} else {
				fmt.Println(dir, "Error:", dlErr.Error())
				return dlErr
			}
		}

		// Find link to previous comic
		links := prevLinkMatch.FindStringSubmatch(s)
		if len(links) < 1 {
			fmt.Println(dir, "No previous URL found")
			return nil
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

func GenericCustomStart(startURL string, startMatch *regexp.Regexp, dir string, fileMatch *regexp.Regexp, filePrefix string, prevLinkMatch *regexp.Regexp, skipExisting bool) error {
	// Load start page HTML
	resp, err := http.Get(startURL)
	if err != nil {
		fmt.Println(dir, "Failed to load page:", err)
		return err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	s := buf.String()

	// Find comic start page
	start := startMatch.FindStringSubmatch(s)
	if len(start) < 1 {
		fmt.Println(dir, "No start URL found")
		return errors.New("no start url found")
	}

	// Start comic download
	return Generic(start[1], dir, fileMatch, filePrefix, prevLinkMatch, skipExisting)
}
