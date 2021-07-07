package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/Alanaktion/comic-archiver/archivers"
)

func main() {
	var all bool
	var skipExisting bool
	flag.BoolVar(&all, "all", false, "Download all supported comics")
	flag.BoolVar(&skipExisting, "continue", false, "Continue download, skipping existing files.")
	flag.Parse()
	var args = flag.Args()

	if !all && len(args) == 0 {
		fmt.Println("Usage: archiver [flags] [comics]")
		fmt.Println("Flags: -all -continue")
		comics := make([]string, len(archivers.Comics))
		i := 0
		for k := range archivers.Comics {
			comics[i] = k
			i++
		}
		fmt.Println("Comics:", comics)
		return
	}

	var comics []string
	if all {
		comics = make([]string, len(archivers.Comics))
		i := 0
		for k := range archivers.Comics {
			comics[i] = k
			i++
		}
	} else {
		comics = args
	}

	var wg sync.WaitGroup
	wg.Add(len(comics))
	for _, c := range comics {
		go archivers.Archive(c, archivers.Comics[c], skipExisting, &wg)
	}
	wg.Wait()
	fmt.Println("Done.")
}
