package archivers

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	neturl "net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

// Generic archiver supporting ComicPress, etc.
func Generic(startURL string, dir string, fileMatch *regexp.Regexp, filePrefix string, prevLinkMatch *regexp.Regexp, skipExisting bool, logger *log.Logger) error {
	os.MkdirAll(dir, os.ModePerm)

	logger.Println(startURL)

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

		// Find comic image
		files := fileMatch.FindStringSubmatch(s)
		if len(files) < 1 {
			logger.Println("No comic image found")
			return errors.New("no comic image found")
		}
		imgUrl := filePrefix + files[1]
		basename := basenameMatch.FindStringSubmatch(files[1])
		path := dir + "/" + basename[1]

		// Download image
		dlErr := downloadFileWait(basename[1], path, imgUrl, 500*time.Millisecond, logger)
		if dlErr != nil {
			if dlErr.Error() == "file exists" {
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
				logger.Println("Error:", dlErr.Error())
				return dlErr
			}
		}

		// Find link to previous comic
		links := prevLinkMatch.FindStringSubmatch(s)
		if len(links) < 1 {
			logger.Println("No previous URL found")
			return nil
		}
		if strings.HasPrefix(links[1], "https://") || strings.HasPrefix(links[1], "http://") {
			url = links[1]
		} else if protocolMatch.MatchString(links[1]) {
			// Malformed URL like "https:/path" (no host) — resolve against current page's origin
			if u, err := neturl.Parse(url); err == nil {
				path := links[1][strings.Index(links[1], ":")+1:]
				url = u.Scheme + "://" + u.Host + path
			} else {
				url = links[1]
			}
		} else {
			// Relative URL: resolve against current page URL
			if base, err := neturl.Parse(url); err == nil {
				if ref, err := neturl.Parse(links[1]); err == nil {
					url = base.ResolveReference(ref).String()
				} else {
					url = startURL + links[1]
				}
			} else {
				url = startURL + links[1]
			}
		}
		recordLastPage(dir, url, logger)

		// Wait a bit
		time.Sleep(500 * time.Millisecond)
	}
}

func GenericCustomStart(startURL string, startMatch *regexp.Regexp, dir string, fileMatch *regexp.Regexp, filePrefix string, prevLinkMatch *regexp.Regexp, skipExisting bool, logger *log.Logger) error {
	// Load start page HTML
	resp, err := http.Get(startURL)
	if err != nil {
		logger.Println("Failed to load page:", err)
		return err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	s := buf.String()

	// Find comic start page
	start := startMatch.FindStringSubmatch(s)
	if len(start) < 1 {
		logger.Println("No start URL found")
		return errors.New("no start url found")
	}

	// Start comic download
	return Generic(start[1], dir, fileMatch, filePrefix, prevLinkMatch, skipExisting, logger)
}
