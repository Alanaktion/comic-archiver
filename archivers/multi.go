package archivers

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
)

// MultiImageGeneric archiver, based on Generic but supporting multiple images per page
func MultiImageGeneric(startURL string, dir string, fileMatch *regexp.Regexp, filePrefix string, prevLinkMatch *regexp.Regexp, namePrefixMatch *regexp.Regexp, skipExisting bool, logger *log.Logger) error {
	os.MkdirAll(dir, os.ModePerm)

	url := startURL
	seen := map[string]bool{url: true}
	last := lastDownloadedPage(dir)
	for {
		// Load page HTML
		resp, err := http.Get(url)
		if err != nil {
			logger.Println("Failed to load page:", err)
			return err
		}
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		s := buf.String()

		// Find and download comic images
		files := fileMatch.FindAllStringSubmatch(s, -1)
		pagePrefix := namePrefix(namePrefixMatch, s)
		for i := range files {
			basename := basenameMatch.FindStringSubmatch(files[i][1])
			name := pagePrefix + basename[1]
			path := dir + "/" + name
			imgurl := filePrefix + files[i][1]
			err := downloadFile(name, path, imgurl, logger)
			if err != nil {
				if err.Error() == "file exists" {
					if !skipExisting {
						logger.Println("File exists:", path)
						return nil
					}
					if last != "" {
						logger.Println("Skipping to URL:", last)
						url = last
						last = ""
						continue
					}
				} else {
					logger.Println("Error:", err.Error())
					return err
				}
			}
		}

		// Find link to previous comic
		links := prevLinkMatch.FindStringSubmatch(s)
		if len(links) < 1 {
			logger.Println("No previous URL found")
			return nil
		}
		if protocolMatch.MatchString(links[1]) {
			url = links[1]
		} else {
			url = startURL + links[1]
		}
		if seen[url] {
			logger.Println("Already visited, stopping:", url)
			return nil
		}
		seen[url] = true
		recordLastPage(dir, url, logger)

		// Wait a bit
		time.Sleep(500 * time.Millisecond)
	}
}
