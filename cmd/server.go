package cmd

import (
	"encoding/json"
	"fmt"
	"inventory/api/providers/database_provider"
	"inventory/pkg/models"
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

func listItemsHandler(w http.ResponseWriter, r *http.Request) {
	provider, err := database_provider.NewMySQLProvider(&config.Database)
	if err != nil {
		http.Error(w, "Failed to obtain provider", http.StatusInternalServerError)
	}
	items, err := provider.SortedListForOutput()
	if err != nil {
		http.Error(w, "Failed to obtain items from provider", http.StatusInternalServerError)
	}
	htmlItems, err := models.OutputItemsToHTML(items)
	if err != nil {
		http.Error(w, "Failed to convert items to strings", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(htmlItems)
	if err != nil {
		http.Error(w, "Failed to JSON encode items", http.StatusInternalServerError)
	}
}

func server() error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/items", listItemsHandler)

	fs := http.FileServer(http.Dir("./web/public"))
	mux.Handle("/", fs)

	fmt.Println("Server listening on http://localhost")
	log.Fatal(http.ListenAndServe(":80", mux))

	return nil
}
