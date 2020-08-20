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
func Generic(startURL string, dir string, fileMatch *regexp.Regexp, filePrefix string, prevLinkMatch *regexp.Regexp) {
	os.MkdirAll("comics/"+dir, os.ModePerm)

	url := startURL
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
		files := fileMatch.FindStringSubmatch(s)
		imgurl := filePrefix + files[1]
		basename := basenameMatch.FindStringSubmatch(files[1])
		path := "comics/" + dir + "/" + basename[1]

		// Download image
		dlErr := downloadFileWait(basename[1], path, imgurl, 500*time.Millisecond)
		if dlErr != nil && dlErr.Error() == "file exists" {
			fmt.Println("File exists:", path)
			return
		}

		// Find link to previous comic
		links := prevLinkMatch.FindStringSubmatch(s)
		if len(links) < 1 {
			fmt.Println("No previous URL found")
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
