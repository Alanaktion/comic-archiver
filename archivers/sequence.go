package archivers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// Sequential archiver downloads images sequentially with predefined prefixes/suffixes
func Sequential(dir string, filePrefix string, pattern string, start int, end int) {
	os.MkdirAll("comics/"+dir, os.ModePerm)

	for i := start; i <= end; i++ {
		name := fmt.Sprintf(pattern, i)
		path := "comics/" + dir + "/" + name
		imgurl := filePrefix + name
		if _, err := os.Stat(path); os.IsNotExist(err) {
			// Image does not exist locally
			fmt.Println("Downloading:", name)
			imgresp, err := http.Get(imgurl)
			if err != nil {
				fmt.Println("Failed to load image:", err)
				return
			}

			buf := new(bytes.Buffer)
			buf.ReadFrom(imgresp.Body)
			err = ioutil.WriteFile(path, buf.Bytes(), 0644)
			if err != nil {
				fmt.Println("Failed to write image:", err)
				return
			}

			// Wait a bit
			time.Sleep(500 * time.Millisecond)
		}
	}
}
