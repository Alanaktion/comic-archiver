package server

import (
	_ "embed"
	"net/http"
)

//go:embed html/app.min.css
var css []byte

// Serve the app CSS file
func cssHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	w.Write(css)
}
