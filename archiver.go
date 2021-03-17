package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/Alanaktion/comic-archiver/archivers"
)

func main() {
	action := "help"
	if len(os.Args) >= 2 {
		action = strings.TrimLeft(os.Args[1], "-")
	}

	if action == "help" || action == "h" {
		fmt.Println("Usage: archiver <action> [params]")
		fmt.Println("Actions: archive, all, help")
		comics := make([]string, len(archivers.Comics))
		i := 0
		for k := range archivers.Comics {
			comics[i] = k
			i++
		}
		fmt.Println("Archivers:", comics)
		return
	}

	if action == "archive" || action == "a" || action == "all" {
		var comics []string
		if action == "all" {
			comics = make([]string, len(archivers.Comics))
			i := 0
			for k := range archivers.Comics {
				comics[i] = k
				i++
			}
		} else {
			comics = os.Args[2:]
		}

		var wg sync.WaitGroup
		wg.Add(len(comics))
		for _, c := range comics {
			go archivers.Archive(c, archivers.Comics[c], &wg)
		}
		wg.Wait()
		fmt.Println("Done.")
		return
	}
}
