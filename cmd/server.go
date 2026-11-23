package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"inventory/api/providers/database_provider"
	"inventory/pkg/models"
	"log"
	"net/http"
	"strings"
	"net/url"
	"gorm.io/gorm"

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


type bulkUpdateRequest struct {
	IDs []int `json:"ids"`
	Field string `json:"field"`
	Value string `json:"value"`
}

// Bulk editing deliberately updates only the requested field, preserving all other item data.
// All IDs are validated before the first write.
func bulkUpdateHandler(provider database_provider.DatabaseProviderInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req bulkUpdateRequest
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		if len(req.IDs) == 0 || len(req.IDs) > 1000 {
			http.Error(w, "Select between 1 and 1000 items", http.StatusBadRequest)
			return
		}
		value := strings.TrimSpace(req.Value)
		maxLen := 100
		switch req.Field {
		case "room", "purchase_source":
		case "info_complete":
			if value != "true" && value != "false" {
				http.Error(w, "Info Complete must be true or false", http.StatusBadRequest)
				return
			}
		case "category", "condition":
			maxLen = 20
		case "website_url", "manual_url":
			maxLen = 2048
			if value != "" {
				u, err := url.Parse(value)
				if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
					http.Error(w, "URL must be an http(s) address", http.StatusBadRequest)
					return
				}
			}
		default:
			http.Error(w, "Unsupported bulk field", http.StatusBadRequest)
			return
		}
		if len(value) > maxLen {
			http.Error(w, "Value is too long", http.StatusBadRequest)
			return
		}
		seen := make(map[int]bool)
        for _, id := range req.IDs {
            if id <= 0 || seen[id] {
                http.Error(w, "Invalid or duplicate item ID", http.StatusBadRequest)
                return
            }
            seen[id] = true
        }
        var updateValue any = value
        if req.Field == "info_complete" { updateValue = value == "true" }
        if err := provider.BulkUpdate(req.IDs, req.Field, updateValue); err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                http.Error(w, "Item not found; no changes made", http.StatusNotFound)
            } else {
                log.Printf("Bulk update failed: %v", err)
                http.Error(w, "Bulk update failed; no changes made", http.StatusInternalServerError)
            }
            return
        }
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]int{"updated": len(req.IDs)})
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
	mux.HandleFunc("POST /api/items/bulk", bulkUpdateHandler(provider))
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
