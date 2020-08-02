package main

import (
	"fmt"
	"os"
	"strings"

	"git.phpizza.com/alan/comic-archiver/archivers"
)

func main() {
	action := "help"
	if len(os.Args) >= 2 {
		action = strings.TrimLeft(os.Args[1], "-")
	}

	if action == "help" || action == "h" {
		fmt.Println("Usage: archiver <action> [params]")
		fmt.Println("Actions: archive, help")
		comics := make([]string, len(archivers.Comics))
		i := 0
		for k := range archivers.Comics {
			comics[i] = k
			i++
		}
		fmt.Println("Archivers: ", comics)
		return
	}

	if action == "archive" || action == "a" {
		comics := os.Args[2:]
		for _, c := range comics {
			fmt.Println("Comic: ", c)
			comic := archivers.Comics[c]
			if comic.Archiver == "ComicPress" {
				archivers.ComicPress(comic.StartURL, c, comic.FileMatch, comic.FilePrefix, comic.PrevLinkMatch)
			}
			if comic.Archiver == "Xkcd" {
				archivers.Xkcd(comic.StartURL, c, comic.FileMatch, comic.FilePrefix, comic.PrevLinkMatch)
			}
		}
		return
	}
}
