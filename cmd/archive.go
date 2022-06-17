package cmd

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/Alanaktion/comic-archiver/archivers"
	"github.com/spf13/cobra"
)

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
		file, err := os.OpenFile("archiver.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(file)

		// Validate args
		if !allFlag && len(args) == 0 {
			fmt.Println("Specify at least one comic to download, or use --all.")
			fmt.Printf("Use '%s list' to see a list of supported comics.\n", os.Args[0])
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
		var lastErr error
		for _, c := range comics {
			val, ok := archivers.Comics[c]
			if !ok {
				fmt.Println("Unknown comic:", c)
				continue
			}
			chanErr := make(chan error)
			go archivers.Archive(c, val, continueFlag, &wg, chanErr)
			err := <-chanErr
			if err != nil {
				lastErr = err
			}
		}
		wg.Wait()
		if lastErr != nil {
			os.Exit(1)
		}
		fmt.Println("Done.")
	},
}

func init() {
	rootCmd.AddCommand(archiveCmd)

	archiveCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "Download all supported comics")
	archiveCmd.Flags().BoolVarP(&continueFlag, "continue", "c", false, "Continue partial downloads, skipping existing files")
}
