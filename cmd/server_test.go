package cmd

import (
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
	return f.listedOutputItems, nil
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
	provider := &fakeProvider{}
	body := strings.NewReader(`{
        "name": "Coffee Maker"
    }`)

	req := httptest.NewRequest(http.MethodPost, "/api/save", body)
	recorder := httptest.NewRecorder()
	handler := saveItemHandler(provider)
	handler(recorder, req)

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
				Name:         "Coffee Maker",
				Quantity:     1,
				PurchaseDate: &time1,
			}, {
				Name:         "Blender",
				Quantity:     2,
				PurchaseDate: &time2,
			},
		},
	}
	want := `{
		{
			"name': "Coffee Maker",
			"quantity": 1,
			"purchaseDate": "01/01/2026"
		},{
			"name': "Blender",
			"quantity": 2,
			"purchaseDate": "02/02/2026"
		}
	}`

	//body := strings.NewReader()

	req := httptest.NewRequest(http.MethodGet, "/api/items", nil)
	recorder := httptest.NewRecorder()
	handler := listItemsHandler(provider)
	handler(recorder, req)

	// Check the HTTP status.
	if recorder.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", recorder.Code, http.StatusOK)
	}

	if reflect.DeepEqual(recorder.Body.String(), want) {
		t.Errorf("got body %q, want %q", recorder.Body.String(), want)
	}
}
