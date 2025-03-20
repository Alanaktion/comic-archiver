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
func MultiImageGeneric(startURL string, dir string, fileMatch *regexp.Regexp, filePrefix string, prevLinkMatch *regexp.Regexp, skipExisting bool, logger *log.Logger) error {
	os.MkdirAll(dir, os.ModePerm)

	url := startURL
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
		files := fileMatch.FindAllStringSubmatch(s, 2)
		for i := range files {
			basename := basenameMatch.FindStringSubmatch(files[i][1])
			path := dir + "/" + basename[1]
			imgurl := filePrefix + files[i][1]
			err := downloadFile(files[i][1], path, imgurl, logger)
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
		recordLastPage(dir, url, logger)

		// Wait a bit
		time.Sleep(500 * time.Millisecond)
	}
}
