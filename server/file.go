package server

import (
	"log"
	"net/http"
	"os"
	"regexp"
)

var validFilePath = regexp.MustCompile(`^/file/(\w+/[\w\. @-]+)$`)

func fileHandler(w http.ResponseWriter, req *http.Request) {
	m := validFilePath.FindStringSubmatch(req.URL.Path)
	if m == nil {
		log.Println("Invalid file path", req.URL.Path)
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	bytes, err := os.ReadFile(m[1])
	if err != nil {
		log.Println(err)
		http.NotFound(w, req)
	} else {
		mimeType := http.DetectContentType(bytes)
		w.Header().Set("Content-Type", mimeType)
		w.Write(bytes)
	}
}
