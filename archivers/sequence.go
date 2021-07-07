package archivers

import (
	"fmt"
	"os"
	"time"
)

// Sequential archiver downloads images sequentially with predefined prefixes/suffixes
func Sequential(dir string, filePrefix string, pattern string, start int, end int, skipExisting bool) {
	os.MkdirAll("comics/"+dir, os.ModePerm)

	for i := start; i <= end; i++ {
		name := fmt.Sprintf(pattern, i)
		path := "comics/" + dir + "/" + name
		imgUrl := filePrefix + name
		err := downloadFileWait(name, path, imgUrl, 500*time.Millisecond)
		if err != nil && err.Error() == "file exists" && !skipExisting {
			fmt.Println("File exists:", path)
			return
		}
		if err != nil {
			fmt.Println("Error:", err.Error())
			return
		}
	}
}
