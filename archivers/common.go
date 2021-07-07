package archivers

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// Download a file
func downloadFile(filename string, path string, url string) error {
	if _, err := os.Stat(path); err == nil {
		return errors.New("file exists")
	}

	fmt.Println("Downloading:", filename)
	imgResp, err := http.Get(url)
	if err != nil {
		fmt.Println("Failed to load image:", err)
		return errors.New("http error")
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(imgResp.Body)
	err = ioutil.WriteFile(path, buf.Bytes(), 0644)
	if err != nil {
		fmt.Println("Failed to write image:", err)
		return errors.New("io error")
	}

	return nil
}

// Download a file and wait after a successful transfer
func downloadFileWait(filename string, path string, url string, delay time.Duration) error {
	err := downloadFile(filename, path, url)
	if err != nil {
		return err
	}

	time.Sleep(delay)
	return nil
}
