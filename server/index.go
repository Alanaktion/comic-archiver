package server

import (
	_ "embed"
	"html/template"
	"net/http"
	"os"

	"github.com/Alanaktion/comic-archiver/archivers"
)

//go:embed html/index.html
var indexHtml []byte
var indexTemplate = template.Must(template.New("index").Parse(string(indexHtml)))

// List all available comics
func indexHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}

	var comics []string
	for key := range archivers.Comics {
		if _, err := os.Stat("comics/" + key); err == nil {
			comics = append(comics, key)
		}
	}

	p := &Page{Comics: comics}
	indexTemplate.Execute(w, p)
}
