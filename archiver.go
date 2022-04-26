package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"sync"

	"github.com/Alanaktion/comic-archiver/archivers"
)

func main() {
	var all bool
	var skipExisting bool
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "Usage:\n")
		fmt.Fprintf(w, "  %s [options] <comic ...>\n", os.Args[0])

		fmt.Fprintf(w, "Options:\n")
		flag.PrintDefaults()

		fmt.Fprintf(w, "Supported comics:\n ")
		for c := range archivers.Comics {
			fmt.Fprintf(w, " %s", c)
		}
		fmt.Fprintf(w, "\n")
	}
	flag.BoolVar(&all, "all", false, "Download all supported comics")
	flag.BoolVar(&skipExisting, "continue", false, "Continue download, skipping existing files.")
	flag.Parse()
	var args = flag.Args()

	file, err := os.OpenFile("archiver.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)

	if !all && len(args) == 0 {
		fmt.Println("Usage: archiver [flags] [comics]")
		flag.PrintDefaults()
		comics := make([]string, len(archivers.Comics))
		i := 0
		for k := range archivers.Comics {
			comics[i] = k
			i++
		}
		sort.Strings(comics)
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
		val, ok := archivers.Comics[c]
		if !ok {
			fmt.Println("Unknown comic:", c)
			continue
		}
		go archivers.Archive(c, val, skipExisting, &wg)
	}
	wg.Wait()
	fmt.Println("Done.")
}
