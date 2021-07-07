package archivers

// Custom one-off archivers that can't be generalized should be implemented here

import (
	"fmt"
	"os"
	"strconv"
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

func intInArray(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
