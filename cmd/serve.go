package cmd

import (
	"fmt"

	"github.com/Alanaktion/comic-archiver/server"
	"github.com/spf13/cobra"
)

var port int

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start a web server to serve comic archives",
	Long:  `This will serve any locally-available comic images from the archive in a simple navigable web interface.`,
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println(fmt.Sprintf("Starting server at http://localhost:%d/", port))
		server.Start(port)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().IntVarP(&port, "port", "p", 8000, "port to listen on")
}
