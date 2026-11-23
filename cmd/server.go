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

func listItemsHandler(
	provider database_provider.DatabaseProviderInterface,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := provider.SortedListForOutput()
		if err != nil {
			http.Error(w, "Failed to obtain items from provider", http.StatusInternalServerError)
			return
		}
		htmlItems := models.OutputItemsToHTML(items)

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		err = json.NewEncoder(w).Encode(htmlItems)
		if err != nil {
			http.Error(w, "Failed to JSON encode items", http.StatusInternalServerError)
			return
		}
	}
}

func saveItemHandler(
	provider database_provider.DatabaseProviderInterface,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var item models.Item
		err := json.NewDecoder(r.Body).Decode(&item)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		_, err = provider.Save(&item)
		if errors.Is(err, models.ErrAlreadyExists) {
			http.Error(w, "item already exists", http.StatusConflict)
			return
		}
		if err != nil {
			log.Printf("Failed to save item: %v", err)
			http.Error(w, "save failed", http.StatusInternalServerError)
			return
		}
	}
}

func deleteItemHandler(
	provider database_provider.DatabaseProviderInterface,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var item models.Item
		err := json.NewDecoder(r.Body).Decode(&item)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		err = provider.Delete(item)
		if err != nil {
			http.Error(w, "delete failed", http.StatusInternalServerError)
			return
		}
	}
}

func server() error {
	provider, err := database_provider.NewMySQLProvider(&config.Database)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/items", listItemsHandler(provider))
	mux.HandleFunc("POST /api/save", saveItemHandler(provider))
	mux.HandleFunc("DELETE /api/delete", deleteItemHandler(provider))
	mux.HandleFunc("GET /api/items/{id}/photos", listPhotosHandler(provider))
	mux.HandleFunc("POST /api/items/{id}/photos", uploadPhotoHandler(provider))
	mux.HandleFunc("GET /api/items/{id}/photos/{filename}", servePhotoHandler(provider))
	mux.HandleFunc("DELETE /api/items/{id}/photos/{photoID}", deletePhotoHandler(provider))

	fs := http.FileServer(http.Dir("./web/public"))
	mux.Handle("/", fs)

	fmt.Println("Server listening on http://localhost")
	log.Fatal(http.ListenAndServe(":80", mux))

	return nil
}
