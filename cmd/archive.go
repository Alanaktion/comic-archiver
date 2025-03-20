package cmd

import (
	"log"
	"os"
	"sync"

	"github.com/Alanaktion/comic-archiver/archivers"
	"github.com/spf13/cobra"
)

var logFlag string
var allFlag bool
var continueFlag bool

// archiveCmd represents the archive command
var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Download the specified comics",
	Long:  `This will download the specified comics from their official sources, storing the images in ./comics`,
	Args:  cobra.OnlyValidArgs,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		comics := []string{}
		for key := range archivers.Comics {
			comics = append(comics, key)
		}
		return comics, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {

		// Initialize logger
		if logFlag != "" {
			file, err := os.OpenFile(logFlag, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatal(err)
			}
			log.SetOutput(file)
			flags := log.Ldate | log.Ltime | log.Lmsgprefix
			log.SetFlags(flags)
		} else {
			flags := log.Ltime | log.Lmsgprefix
			log.SetFlags(flags)
		}

		// Validate args
		if !allFlag && len(args) == 0 {
			log.Println("Specify at least one comic to download, or use --all.")
			log.Printf("Use '%s list' to see a list of supported comics.\n", os.Args[0])
			os.Exit(1)
			return
		}

		// Build comic list
		var comics []string
		if allFlag {
			comics = make([]string, len(archivers.Comics))
			i := 0
			for k := range archivers.Comics {
				comics[i] = k
				i++
			}
		} else {
			comics = args
		}

		// Start archive workers
		var wg sync.WaitGroup
		wg.Add(len(comics))
		for _, c := range comics {
			val, ok := archivers.Comics[c]
			if !ok {
				log.Println("Unknown comic:", c)
				continue
			}
			go archivers.Archive(c, val, continueFlag, &wg)
		}
		wg.Wait()
		log.Println("Done.")
	},
}

func init() {
	rootCmd.AddCommand(archiveCmd)

	archiveCmd.Flags().StringVarP(&logFlag, "log", "l", "", "Log file to output to, otherwise stderr")
	archiveCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "Download all supported comics")
	archiveCmd.Flags().BoolVarP(&continueFlag, "continue", "c", false, "Continue partial downloads, skipping existing files")
}
