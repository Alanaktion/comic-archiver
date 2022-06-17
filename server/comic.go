package server

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/Alanaktion/comic-archiver/archivers"
)

//go:embed html/comic.html
var comicHtml []byte
var comicTemplate = template.Must(template.New("comic").Parse(string(comicHtml)))

//go:embed html/comic_page.html
var comicPageHtml []byte
var comicPageTemplate = template.Must(template.New("comic_page").Parse(string(comicPageHtml)))

func loadComic(title string) (*Page, error) {
	comic, ok := archivers.Comics[title]
	if !ok {
		return nil, os.ErrInvalid
	}
	path := "comics/" + title
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, os.ErrInvalid
	}
	files, err := os.ReadDir("comics/" + title)
	if err != nil {
		return nil, err
	}
	pages := make([]string, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		pages = append(pages, file.Name())
	}
	return &Page{Title: title, Comic: comic, Pages: pages}, nil
}

var validComicPath = regexp.MustCompile(`^/comic/(\w+)(/(\d+))?$`)

func comicHandler(w http.ResponseWriter, req *http.Request) {
	m := validComicPath.FindStringSubmatch(req.URL.Path)
	if m == nil {
		fmt.Println(404, req.URL.Path, "pattern match failed")
		http.NotFound(w, req)
		return
	}

	title := m[1]
	p, err := loadComic(title)
	if err != nil {
		fmt.Println(404, req.URL.Path, "comic not found")
		http.NotFound(w, req)
		return
	}

	// Main comic file listing
	if m[2] == "" || m[2] == "/" {
		comicTemplate.Execute(w, p)
		return
	}

	// Specific comic image page
	index, err := strconv.Atoi(m[3])
	if err != nil {
		http.NotFound(w, req)
		return
	}
	if index < 0 || index >= len(p.Pages) {
		fmt.Println(404, req.URL.Path, "page not in comic")
		http.NotFound(w, req)
		return
	}
	p.Index = index
	p.Page = p.Pages[index]
	if index > 0 {
		p.PrevPage = p.Pages[index-1]
		p.PrevPageIndex = index - 1
	}
	if index < len(p.Pages)-1 {
		p.NextPage = p.Pages[index+1]
		p.NextPageIndex = index + 1
	}

	comicPageTemplate.Execute(w, p)
}
