package server

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"

	"github.com/Alanaktion/comic-archiver/archivers"
)

type Page struct {
	Title         string
	Comics        []string
	Comic         archivers.Comic
	Pages         []string
	Index         int
	Page          string
	PrevPage      string
	NextPage      string
	PrevPageIndex int
	NextPageIndex int
}

func Start(port int) {
	// Configure routes
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/app.css", cssHandler)
	http.HandleFunc("/comic/", comicHandler)
	http.HandleFunc("/file/", fileHandler)

	// Start server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
