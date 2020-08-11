package archivers

import (
	"fmt"
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
		err := downloadFileWait(name, path, imgurl, 500*time.Millisecond)
		if err != nil {
			return
		}
	}
}
