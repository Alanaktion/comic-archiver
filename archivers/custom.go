package archivers

// Custom one-off archivers that can't be generalized should be implemented here

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// AliceGrove archives Jeph's custom-coded, semi-broken site
func AliceGrove(dir string, filePrefix string, end int, skipExisting bool, logger *log.Logger) error {
	os.MkdirAll(dir, os.ModePerm)

	jpegs := []int{
		35, 70, 78, 83, 84, 98, 100, 107, 113, 124,
		126, 127, 128, 129, 130, 131, 132, 134, 136,
		141, 145, 153, 159, 164, 168, 169, 170, 171,
		172, 173, 174, 175, 176, 177, 178, 179, 180,
		181, 182, 183, 186, 196,
	}

	for i := 1; i <= end; i++ {
		// 109 and 165 are unique, 137 doesn't exist :P
		if i == 109 || i == 165 || i == 137 {
			continue
		}

		var name string
		if intInArray(i, jpegs) {
			name = strconv.FormatInt(int64(i), 10) + ".jpg"
		} else {
			name = strconv.FormatInt(int64(i), 10) + ".png"
		}
		path := dir + "/" + name
		imgUrl := filePrefix + name
		err := downloadFileWait(name, path, imgUrl, 500*time.Millisecond, logger)
		if err != nil {
			if err.Error() == "file exists" {
				if !skipExisting {
					logger.Println("File exists:", path)
					return nil
				}
			} else {
				logger.Println("Error:", err.Error())
				return err
			}
		}
	}

	// Handle non-standard images
	extra := []string{"109-1.jpg", "109-2.png", "165-1.png", "165-2.jpg"}
	for i := range extra {
		name := extra[i]
		path := dir + "/" + name
		imgUrl := filePrefix + name
		err := downloadFileWait(name, path, imgUrl, 500*time.Millisecond, logger)
		if err != nil {
			if err.Error() == "file exists" {
				if !skipExisting {
					logger.Println("File exists:", path)
					return nil
				}
			} else {
				logger.Println("Error:", err.Error())
				return err
			}
		}
	}

	return nil
}

// Floraverse *could* work fine via the Generic downloader, but I want a better way of naming the files rather than the hashes used in the server filenames.
func Floraverse(startURL string, dir string, skipExisting bool, logger *log.Logger) error {
	os.MkdirAll(dir, os.ModePerm)

	fileMatch := regexp.MustCompile(`src="https://floraverse.com/filestore/([^"]+\.(jpg|png|gif))`)
	filePrefix := "https://floraverse.com/filestore/"
	prevLinkMatch := regexp.MustCompile(`href="(https://floraverse.com/comic/[0-9a-zA-Z/_-]+)">◀ previous( by date|<)`)
	namePathMatch := regexp.MustCompile(`page.identifier = "https://floraverse.com/comic/([0-9a-zA-Z/_-]+)"`)

	// May want to just read the archive page and get every link from there to avoid navigation jank.
	// https://floraverse.com/comic/%40%40by-date/

	url := startURL
	seen := map[string]bool{url: true}
	last := lastDownloadedPage(dir)
	for {
		s, err := httpGetAsString(url)
		if err != nil {
			logger.Println("Failed to load page:", err)
			return err
		}

		files := fileMatch.FindStringSubmatch(s)
		if len(files) < 1 {
			logger.Println("No comic image found")
			return nil
		}
		imgUrl := filePrefix + files[1]
		destName := strings.ReplaceAll(strings.Trim(namePathMatch.FindStringSubmatch(s)[1], "/"), "/", "_") + "." + files[2]
		path := dir + "/" + destName

		dlErr := downloadFileWait(destName, path, imgUrl, 500*time.Millisecond, logger)
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
				return nil
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

// Monochrome post URLs aren't in chronological order, so files are prefixed with the post's publish date and time to make them sort in reading order.
func Monochrome(startURL string, dir string, skipExisting bool, logger *log.Logger) error {
	os.MkdirAll(dir, os.ModePerm)

	fileMatch := regexp.MustCompile(`id="comic">\s*(?:<a [^>]*>)?\s*<img src="([^"]+\.(jpg|jpeg|png|gif))"`)
	prevLinkMatch := regexp.MustCompile(`href="([^"]+)" class="comic-nav-base comic-nav-previous"`)
	dateMatch := regexp.MustCompile(`class="post-date">([^<]+)</span>\s*<span class="posted-at">at&nbsp;</span><span class="post-time">([^<]+)</span>`)
	slugMatch := regexp.MustCompile(`/comic/([^/]+)/`)

	url := startURL
	seen := map[string]bool{url: true}
	last := lastDownloadedPage(dir)
	for {
		s, err := httpGetAsString(url)
		if err != nil {
			logger.Println("Failed to load page:", err)
			return err
		}

		files := fileMatch.FindStringSubmatch(s)
		if len(files) < 1 {
			// Some posts are text-only announcements; skip them
			logger.Println("No comic image found, skipping:", url)
			links := prevLinkMatch.FindStringSubmatch(s)
			if len(links) < 1 {
				logger.Println("No previous URL found")
				return nil
			}
			url = links[1]
			if seen[url] {
				logger.Println("Already visited, stopping:", url)
				return nil
			}
			seen[url] = true
			recordLastPage(dir, url, logger)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		dates := dateMatch.FindStringSubmatch(s)
		if len(dates) < 1 {
			logger.Println("No post date found")
			return errors.New("no post date found")
		}
		published, err := time.Parse("January 2, 2006 3:04 pm", dates[1]+" "+dates[2])
		if err != nil {
			logger.Println("Failed to parse post date:", err)
			return err
		}

		// The home page has no slug, so name by the image file instead
		name := published.Format("2006-01-02_1504") + "_" + files[1][strings.LastIndex(files[1], "/")+1:]
		if slug := slugMatch.FindStringSubmatch(url); len(slug) > 1 {
			name = published.Format("2006-01-02_1504") + "_" + slug[1] + "." + files[2]
		}
		path := dir + "/" + name

		dlErr := downloadFileWait(name, path, files[1], 500*time.Millisecond, logger)
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

		links := prevLinkMatch.FindStringSubmatch(s)
		if len(links) < 1 {
			logger.Println("No previous URL found")
			return nil
		}
		url = links[1]
		if seen[url] {
			logger.Println("Already visited, stopping:", url)
			return nil
		}
		seen[url] = true
		recordLastPage(dir, url, logger)

		time.Sleep(500 * time.Millisecond)
	}
}

// LoadingArtist serves resized, hash-named renditions on the page, so build the URL of the full-size original from the comic's slug instead.
func LoadingArtist(startURL string, dir string, skipExisting bool, logger *log.Logger) error {
	os.MkdirAll(dir, os.ModePerm)

	slugMatch := regexp.MustCompile(`content="https://loadingartist.com/comic/([^/"]+)/"`)
	dateMatch := regexp.MustCompile(`<time datetime=(\d{4}-\d{2}-\d{2})>`)
	prevLinkMatch := regexp.MustCompile(`<a title=Previous href=(/comic/[^ >]+/)`)

	url := startURL
	seen := map[string]bool{url: true}
	last := lastDownloadedPage(dir)
	for {
		s, err := httpGetAsString(url)
		if err != nil {
			logger.Println("Failed to load page:", err)
			return err
		}

		slug := slugMatch.FindStringSubmatch(s)
		date := dateMatch.FindStringSubmatch(s)
		if len(slug) < 2 || len(date) < 2 {
			logger.Println("No comic slug or date found:", url)
			return errors.New("no comic found")
		}

		// Originals are not linked anywhere, so probe for the extension
		var imgUrl, ext string
		for _, e := range []string{"png", "jpg", "gif"} {
			u := "https://loadingartist.com/comic/" + slug[1] + "/" + slug[1] + "." + e
			if resp, err := http.Head(u); err == nil {
				resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
					imgUrl, ext = u, e
					break
				}
			}
		}
		if imgUrl == "" {
			logger.Println("No comic image found:", url)
			return errors.New("no comic image found")
		}

		name := date[1] + "_" + slug[1] + "." + ext
		path := dir + "/" + name
		dlErr := downloadFileWait(name, path, imgUrl, 500*time.Millisecond, logger)
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

		links := prevLinkMatch.FindStringSubmatch(s)
		if len(links) < 1 {
			logger.Println("No previous URL found")
			return nil
		}
		url = "https://loadingartist.com" + links[1]
		if seen[url] {
			logger.Println("Already visited, stopping:", url)
			return nil
		}
		seen[url] = true
		recordLastPage(dir, url, logger)

		time.Sleep(500 * time.Millisecond)
	}
}

func intInArray(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func httpGetAsString(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	s := buf.String()
	return s, nil
}
