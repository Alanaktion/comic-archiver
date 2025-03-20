package archivers

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"os"
	"time"
)

// Download a file
func downloadFile(filename string, path string, url string, logger *log.Logger) error {
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
