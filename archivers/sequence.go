package archivers

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Sequential archiver downloads images sequentially with predefined prefixes/suffixes
func Sequential(dir string, filePrefix string, pattern string, start int, end int, skipExisting bool, logger *log.Logger) error {
	os.MkdirAll(dir, os.ModePerm)

	for i := start; i <= end; i++ {
		name := fmt.Sprintf(pattern, i)
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
