package archivers

// Custom one-off archivers that can't be generalized should be implemented here

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// AliceGrove archives Jeph's custom-coded, semi-broken site
func AliceGrove(dir string, filePrefix string, end int, skipExisting bool) {
	os.MkdirAll("comics/"+dir, os.ModePerm)

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
		path := "comics/" + dir + "/" + name
		imgUrl := filePrefix + name
		err := downloadFileWait(name, path, imgUrl, 500*time.Millisecond)
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

	// Handle non-standard images
	extra := []string{"109-1.jpg", "109-2.png", "165-1.png", "165-2.jpg"}
	for i := range extra {
		name := extra[i]
		path := "comics/" + dir + "/" + name
		imgUrl := filePrefix + name
		err := downloadFileWait(name, path, imgUrl, 500*time.Millisecond)
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
}

// Floraverse *could* work fine via the Generic downloader, but I want a better way of naming the files rather than the hashes used in the server filenames.
func Floraverse(startURL string, dir string, skipExisting bool) {
	os.MkdirAll("comics/"+dir, os.ModePerm)

	fileMatch := regexp.MustCompile(`src="https://floraverse.com/filestore/([^"]+\.(jpg|png|gif))`)
	filePrefix := "https://floraverse.com/filestore/"
	prevLinkMatch := regexp.MustCompile(`href="(https://floraverse.com/comic/[0-9a-zA-Z/_-]+)">◀ previous by date`)
	namePathMatch := regexp.MustCompile(`page.identifier = "https://floraverse.com/comic/([0-9a-zA-Z/_-]+)"`)

	url := startURL
	for {
		s, err := httpGetAsString(url)
		if err != nil {
			fmt.Println(dir, "Failed to load page:", err)
			os.Exit(1)
		}

		files := fileMatch.FindStringSubmatch(s)
		if len(files) < 1 {
			fmt.Println(dir, "No comic image found")
			return
		}
		imgUrl := filePrefix + files[1]
		destName := strings.ReplaceAll(strings.Trim(namePathMatch.FindStringSubmatch(s)[1], "/"), "/", "_") + "." + files[2]
		path := "comics/" + dir + "/" + destName

		dlErr := downloadFileWait(destName, path, imgUrl, 500*time.Millisecond)
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
