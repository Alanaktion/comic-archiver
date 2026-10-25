package archivers

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

// Some sites (e.g. Monochrome) reject Go's default User-Agent, so identify ourselves instead.
type userAgentTransport struct{ base http.RoundTripper }

func (t userAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	req.Header.Set("User-Agent", "comic-archiver")
	return t.base.RoundTrip(req)
}

func init() {
	http.DefaultClient.Transport = userAgentTransport{http.DefaultTransport}
	// Never let a stalled server block a comic (and the whole run) forever
	http.DefaultClient.Timeout = 2 * time.Minute
}

// Download a file
func downloadFile(filename string, path string, url string, logger *log.Logger) error {
	// TODO: support glob matching so any file extension can be considered,
	// e.g. when used with `python3 -m imgfind.teeny -f webp {}`.
	// Alternatively, start using my gallery-dl extractor + archive matching.
	if _, err := os.Stat(path); err == nil {
		return errors.New("file exists")
	}

	logger.Println("Downloading:", filename)
	imgResp, err := http.Get(url)
	if err != nil {
		logger.Println("Failed to load image:", err)
		return errors.New("http error")
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(imgResp.Body)
	err = os.WriteFile(path, buf.Bytes(), 0644)
	if err != nil {
		logger.Println("Failed to write image:", err)
		return errors.New("io error")
	}

	return nil
}

// Download a file and wait after a successful transfer
func downloadFileWait(filename string, path string, url string, delay time.Duration, logger *log.Logger) error {
	err := downloadFile(filename, path, url, logger)
	if err != nil {
		return err
	}

	time.Sleep(delay)
	return nil
}

// Check for last known URL for the given comic
func lastDownloadedPage(comic string) string {
	data, err := os.ReadFile(comic + "/.last_url")
	if err != nil {
		return ""
	}
	return string(data)
}

// Record last known URL for the given comic
func recordLastPage(comic string, url string, logger *log.Logger) error {
	err := os.WriteFile(comic+"/.last_url", []byte(url), 0644)
	if err != nil {
		logger.Println("Failed to write last URL:", err)
	}
	return err
}

// Build a filename prefix from the page HTML so files sort in reading order
// even when the site's own filenames don't. Numeric prefixes are zero-padded.
func namePrefix(match *regexp.Regexp, page string) string {
	if match == nil {
		return ""
	}
	m := match.FindStringSubmatch(page)
	if len(m) < 2 {
		return ""
	}
	if n, err := strconv.Atoi(m[1]); err == nil {
		return fmt.Sprintf("%06d_", n)
	}
	return m[1] + "_"
}
