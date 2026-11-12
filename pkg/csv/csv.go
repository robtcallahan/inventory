package csv

import (
	stdcsv "encoding/csv"
	"fmt"
	"inventory/pkg/models"
	"os"
	"reflect"
	"strings"

	"github.com/gocarina/gocsv"
)

func ImportCSV(fileName string) ([]models.CSVItem, error) {
	if err := checkCSVHeaders(fileName, models.CSVItem{}); err != nil {
		return nil, err
	}
	csvFilePtr, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("could not open file %s: %w", fileName, err)
	}
	defer csvFilePtr.Close()

	var items []models.CSVItem
	if err := gocsv.UnmarshalFile(csvFilePtr, &items); err != nil {
		return nil, fmt.Errorf("could not unmarshal CSV file %s: %w", fileName, err)
	}
	return items, err
}

// normalizeHeader canonicalizes a CSV header label for matching: surrounding
// whitespace is ignored and matching is case-insensitive.
func normalizeHeader(h string) string {
	return strings.ToLower(strings.TrimSpace(h))
}

// csvHeaders returns the `csv:` struct tags of a gocsv-mapped struct: the
// headers the bank's file must contain.
func csvHeaders(sample any) []string {
	t := reflect.TypeOf(sample)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	var headers []string
	for field := range t.Fields() {
		if tag, ok := field.Tag.Lookup("csv"); ok && tag != "" && tag != "-" {
			headers = append(headers, tag)
		}
	}
	return headers
}

// checkCSVHeaders reads a CSV file's header row and verifies it contains every
// header the given gocsv struct expects. gocsv maps columns by header name
// but stays silent when one is missing (the field just keeps its zero value),
// so a renamed bank column would silently zero a field; this makes it loud.
func checkCSVHeaders(csvFile string, sample any) error {
	f, err := os.Open(csvFile)
	if err != nil {
		return fmt.Errorf("could not open file %s: %w", csvFile, err)
	}
	defer f.Close()
	r := stdcsv.NewReader(f)
	header, err := r.Read()
	if err != nil {
		return fmt.Errorf("could not read header row of %s: %w", csvFile, err)
	}
	have := make(map[string]bool)
	for _, h := range header {
		have[normalizeHeader(h)] = true
	}
	var missing []string
	for _, want := range csvHeaders(sample) {
		if !have[normalizeHeader(want)] {
			missing = append(missing, want)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("CSV file %s is missing expected headers: %s", csvFile, strings.Join(missing, ", "))
	}
	return nil
}
