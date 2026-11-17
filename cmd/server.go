package cmd

import (
	"encoding/json"
	"errors"
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

func saveItemHandler(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	provider, err := database_provider.NewMySQLProvider(&config.Database)
	if err != nil {
		http.Error(w, "Failed to obtain provider", http.StatusInternalServerError)
	}

	err = provider.Save(item)
	if errors.Is(err, models.ErrAlreadyExists) {
		http.Error(w, "item already exists", http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, "save failed", http.StatusInternalServerError)
		return
	}
}

func server() error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/items", listItemsHandler)
	mux.HandleFunc("POST /api/save", saveItemHandler)

	fs := http.FileServer(http.Dir("./web/public"))
	mux.Handle("/", fs)

	fmt.Println("Server listening on http://localhost")
	log.Fatal(http.ListenAndServe(":80", mux))

	return nil
}
