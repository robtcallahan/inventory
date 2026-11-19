package cmd

import (
	"errors"
	"inventory/api/providers/database_provider"
	"inventory/pkg/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeProvider struct {
	savedItem models.Item
	saveErr   error
}

func (f *fakeProvider) Save(item models.Item) error {
	f.savedItem = item
	return f.saveErr
}
func (f *fakeProvider) SaveItems(items []models.Item) ([]models.Item, error) {
	return nil, nil
}

func (f *fakeProvider) List() ([]models.Item, error) {
	return nil, nil
}

func (f *fakeProvider) ListByCategory(category string) ([]models.Item, error) {
	return nil, nil
}

func (f *fakeProvider) TotalByCategory(category string) (int, error) {
	return 0, nil
}

func (f *fakeProvider) SortedListForOutput() ([]models.OutputItem, error) {
	return nil, nil
}

var _ database_provider.DatabaseProviderInterface = (*fakeProvider)(nil)

func TestSaveItemHandler_InvalidJSON(t *testing.T) {
	body := strings.NewReader(`{"name":`)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/save",
		body,
	)

	recorder := httptest.NewRecorder()

	provider := &fakeProvider{}
	handler := saveItemHandler(provider)
	handler(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf(
			"got status %d, want %d",
			recorder.Code,
			http.StatusBadRequest,
		)
	}
}

func TestSaveItemHandler_Success(t *testing.T) {
	provider := &fakeProvider{}

	body := strings.NewReader(`{
        "name": "Coffee Maker"
    }`)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/save",
		body,
	)

	recorder := httptest.NewRecorder()
	handler := saveItemHandler(provider)

	handler(recorder, req)

	// Check the HTTP status.
	if recorder.Code != http.StatusOK {
		t.Errorf(
			"got status %d, want %d",
			recorder.Code,
			http.StatusOK,
		)
	}

	// Check what was passed to the provider.
	if provider.savedItem.Name != "Coffee Maker" {
		t.Errorf(
			"saved name = %q, want %q",
			provider.savedItem.Name,
			"Coffee Maker",
		)
	}
}

func TestSaveItemHandler_Duplicate(t *testing.T) {
	provider := &fakeProvider{
		saveErr: models.ErrAlreadyExists,
	}

	body := strings.NewReader(`{
        "name": "Coffee Maker"
    }`)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/save",
		body,
	)

	recorder := httptest.NewRecorder()
	handler := saveItemHandler(provider)

	handler(recorder, req)
	if recorder.Code != http.StatusConflict {
		t.Errorf("go status %d, want %d", recorder.Code, http.StatusConflict)
	}
}

func TestSaveItemHandler_Error(t *testing.T) {
	provider := &fakeProvider{
		saveErr: errors.New("save failed"),
	}
	body := strings.NewReader(`{
        "name": "Coffee Maker"
    }`)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/save",
		body,
	)

	recorder := httptest.NewRecorder()
	handler := saveItemHandler(provider)

	handler(recorder, req)
	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("go status %d, want %d", recorder.Code, http.StatusInternalServerError)
	}
}
