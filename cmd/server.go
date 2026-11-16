package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Service HTTP web pages",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return server()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().BoolVar(&options.Debug, "debug", false, "Debug mode")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Hello, world!")
	if err != nil {
		log.Fatal(err)
	}
}

func server() error {
	// Specify the local directory to serve
	fs := http.FileServer(http.Dir("./web/public"))

	// Route all root requests to the file server
	http.Handle("/", fs)

	// Start the server on port 8080
	fmt.Println("Server listening on http://localhost:80")
	log.Fatal(http.ListenAndServe(":80", nil))

	return nil
}
