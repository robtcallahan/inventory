package cmd

import (
	"encoding/json"
	"errors"
	"inventory/api/providers/database_provider"
	"inventory/pkg/models"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

type fakeProvider struct {
	savedItem         models.Item
	listedOutputItems []models.OutputItem
	saveCalled        bool
	saveErr           error
	listErr           error
}

func (f *fakeProvider) Save(item models.Item) error {
	f.savedItem = item
	f.saveCalled = true
	return f.saveErr
}
func (f *fakeProvider) SaveItems(items []models.Item) ([]models.Item, error) {
	return nil, nil
}

func (f *fakeProvider) List() ([]models.Item, error) {
	return nil, f.listErr
}

func (f *fakeProvider) ListByCategory(category string) ([]models.Item, error) {
	return nil, nil
}

func (f *fakeProvider) TotalByCategory(category string) (int, error) {
	return 0, nil
}

func (f *fakeProvider) SortedListForOutput() ([]models.OutputItem, error) {
	return f.listedOutputItems, f.listErr
}

var _ database_provider.DatabaseProviderInterface = (*fakeProvider)(nil)

func TestSaveItemHandler_InvalidJSON(t *testing.T) {
	provider := &fakeProvider{}

	body := strings.NewReader(`{"name":`)

	req := httptest.NewRequest(http.MethodPost, "/api/save", body)
	recorder := httptest.NewRecorder()
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

	req := httptest.NewRequest(http.MethodPost, "/api/save", body)
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
		saveErr: errors.New("this thing doesn't work"),
	}
	body := strings.NewReader(`{
        "name": "Coffee Maker"
    }`)

	req := httptest.NewRequest(http.MethodPost, "/api/save", body)
	recorder := httptest.NewRecorder()
	handler := saveItemHandler(provider)
	handler(recorder, req)

	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("go status %d, want %d", recorder.Code, http.StatusInternalServerError)
	}
}

func TestSaveItemHandler_BadRequest(t *testing.T) {
	provider := &fakeProvider{
		saveErr: errors.New("save failed"),
	}
	body := strings.NewReader(`{
        "name": "Coffee Maker",
		"value": "49.99"
    }`)

	req := httptest.NewRequest(http.MethodPost, "/api/save", body)
	recorder := httptest.NewRecorder()
	handler := saveItemHandler(provider)
	handler(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("go status %d, want %d", recorder.Code, http.StatusBadRequest)
	}
	if provider.saveCalled {
		t.Error("save called when it shouldn't")
	}
}

func TestListItemsHandler_Success(t *testing.T) {
	time1, _ := time.Parse("01/02/2006", "01/01/2026")
	time2, _ := time.Parse("01/02/2006", "02/02/2026")
	provider := &fakeProvider{
		listedOutputItems: []models.OutputItem{
			{
				ID:           1,
				Name:         "Coffee Maker",
				Quantity:     1,
				Value:        25,
				TotalValue:   25,
				PurchaseDate: &time1,
			}, {
				ID:           2,
				Name:         "Blender",
				Quantity:     2,
				Value:        10,
				TotalValue:   20,
				PurchaseDate: &time2,
			},
		},
	}
	//	want := `[{"id":1,"number":1,"quantity":1,"name":"Coffee Maker","value":25,"total_value":25,"description":"","manufacturer":"","model":"","category":"","purchase_date":"01/01/2026","serial_number":"","details":""},{"id":2,"number":2,"quantity":2,"name":"Blender","value":10,"total_value":20,"description":"","manufacturer":"","model":"","category":"","purchase_date":"02/02/2026","serial_number":"","details":""}]
	//`
	want := models.OutputItemsToHTML(provider.listedOutputItems)

	req := httptest.NewRequest(http.MethodGet, "/api/items", nil)
	recorder := httptest.NewRecorder()
	handler := listItemsHandler(provider)
	handler(recorder, req)

	// Check the HTTP status.
	if recorder.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", recorder.Code, http.StatusOK)
	}

	var got []models.HTMLItem
	err := json.Unmarshal(recorder.Body.Bytes(), &got)
	if err != nil {
		t.Fatalf("failed to unmarshal response body: %s", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got body %+v, want %+v", got, want)
	}
}

func TestListItemsHandler_Error(t *testing.T) {
	provider := &fakeProvider{
		listErr: errors.New("SortedListForOutput failed"),
	}

	req := httptest.NewRequest(http.MethodGet, "/api/items", nil)
	recorder := httptest.NewRecorder()
	handler := listItemsHandler(provider)
	handler(recorder, req)

	// Check the HTTP status.
	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("got status %d, want %d", recorder.Code, http.StatusInternalServerError)
	}
}
