package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/Alanaktion/comic-archiver/archivers"
	"github.com/spf13/cobra"
)

var local bool
var oneLine bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available comics",
	Long:  `List either all comics that are supported by the archiver, or all comics that are locally available.`,
	Run: func(cmd *cobra.Command, args []string) {
		var comics []string
		if local {
			for key := range archivers.Comics {
				if _, err := os.Stat(key); err == nil {
					comics = append(comics, key)
				}
			}
		} else {
			for key := range archivers.Comics {
				comics = append(comics, key)
			}
		}
		sort.Strings(comics)
		if oneLine {
			for _, comic := range comics {
				fmt.Println(comic)
			}
		} else {
			fmt.Println(strings.Join(comics, ", "))
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVarP(&local, "local", "l", false, "List comics available locally")
	listCmd.Flags().BoolVarP(&oneLine, "one-line", "1", false, "List one comic per line")
}
