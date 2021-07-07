package archivers

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"
)

// MultiImageGeneric archiver, based on Generic but supporting multiple images per page
func MultiImageGeneric(startURL string, dir string, fileMatch *regexp.Regexp, filePrefix string, prevLinkMatch *regexp.Regexp, skipExisting bool) {
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

		// Find and download comic images
		files := fileMatch.FindAllStringSubmatch(s, 2)
		for i := range files {
			basename := basenameMatch.FindStringSubmatch(files[i][1])
			path := "comics/" + dir + "/" + basename[1]
			imgurl := filePrefix + files[i][1]
			err := downloadFile(files[i][1], path, imgurl)
			if err != nil {
				if err.Error() == "file exists" {
					if !skipExisting {
						fmt.Println("File exists:", path)
						return
					}
				} else {
					fmt.Println("Error:", err.Error())
					return
				}
			}
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
