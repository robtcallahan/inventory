package cmd

import (
	"inventory/pkg/models"
	"path/filepath"
	"testing"
)

func Test_ItemsExportToCSVFile(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "inventory.csv")

	got := models.ItemsExportToCSVFile(outputFile, items)
}
