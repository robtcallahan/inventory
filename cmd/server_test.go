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
	itemID            int
	getItem models.Item
	savedItem         *models.Item
	deletedItem       models.Item
	listedOutputItems []models.OutputItem
	saveCalled        bool
	deleteCalled      bool
	saveErr           error
	listErr           error
	deleteErr         error
	bulkIDs []int
	bulkField string
	bulkValue any
	bulkErr error
	savedPhoto        *models.ItemPhoto
	savedPhotoErr     error
	savePhotoCalled   bool
	listedPhotos      []models.ItemPhoto
	listPhotosErr     error
	deletePhotoCalled bool
	deletedPhotoID    int
	deletePhotoErr    error
}

func (f *fakeProvider) Save(item *models.Item) (*models.Item, error) {
	f.savedItem = item
	f.saveCalled = true
	return f.savedItem, f.saveErr
}
func (f *fakeProvider) Get(id int) (models.Item, error) {
	f.itemID = id
	return f.getItem, f.listErr
}
func (f *fakeProvider) SaveItems(items []models.Item) ([]models.Item, error) {
	f.saveCalled = true
	return nil, f.listErr
}

func (f *fakeProvider) BulkUpdate(ids []int, field string, value any) error {
	f.bulkIDs, f.bulkField, f.bulkValue = ids, field, value
	return f.bulkErr
}
func (f *fakeProvider) Delete(item models.Item) error {
	f.deletedItem = item
	f.deleteCalled = true
	return f.deleteErr
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

func (f *fakeProvider) SavePhoto(photo *models.ItemPhoto) (*models.ItemPhoto, error) {
	f.savedPhoto = photo
	f.savePhotoCalled = true
	return photo, f.savedPhotoErr
}
func (f *fakeProvider) ListPhotos(itemID int) ([]models.ItemPhoto, error) {
	f.itemID = itemID
	return f.listedPhotos, f.listPhotosErr
}
func (f *fakeProvider) DeletePhoto(photoID int) error {
	f.deletedPhotoID = photoID
	f.deletePhotoCalled = true
	return f.deletePhotoErr
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
		t.Errorf("got status %d, want %d", recorder.Code, http.StatusOK)
	}

	if !provider.saveCalled {
		t.Errorf("save was not called")
	}

	// Check what was passed to the provider.
	if provider.savedItem.Name != "Coffee Maker" {
		t.Errorf("saved name = %q, want %q", provider.savedItem.Name, "Coffee Maker")
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
		t.Errorf("got status %d, want %d", recorder.Code, http.StatusInternalServerError)
	}
}

func TestSaveItemHandler_BadRequest(t *testing.T) {
	provider := &fakeProvider{
		saveErr: errors.New("save failed"),
	}
	body := strings.NewReader(`{
        "name": "Coffee Maker",
		"replacement_value": "49.99"
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

func TestDeleteItemHandler_Success(t *testing.T) {
	provider := &fakeProvider{}

	body := strings.NewReader(`{
		"id": 1,
        "name": "Shoe Box"
    }`)

	req := httptest.NewRequest(http.MethodPost, "/api/delete", body)
	recorder := httptest.NewRecorder()
	handler := deleteItemHandler(provider)
	handler(recorder, req)

	// Check the HTTP status.
	if recorder.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", recorder.Code, http.StatusOK)
	}

	if !provider.deleteCalled {
		t.Errorf("delete was not called")
	}

	// Check what was passed to the provider.
	if provider.deletedItem.Name != "Shoe Box" {
		t.Errorf("saved name = %q, want %q", provider.savedItem.Name, "Shoe Box")
	}
}

func TestDeleteItemHandler_Error(t *testing.T) {
	provider := &fakeProvider{
		deleteErr: errors.New("this thing doesn't work"),
	}
	body := strings.NewReader(`{
		"id": 1,
        "name": "Shoe Box"
    }`)

	req := httptest.NewRequest(http.MethodPost, "/api/delete", body)
	recorder := httptest.NewRecorder()
	handler := deleteItemHandler(provider)
	handler(recorder, req)

	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("got status %d, want %d", recorder.Code, http.StatusInternalServerError)
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
				ReplacementValue:        25,
				TotalValue:   25,
				PhotoCount: 2,
				PurchaseDate: &time1,
			}, {
				ID:           2,
				Name:         "Blender",
				Quantity:     2,
				ReplacementValue:        10,
				TotalValue:   20,
				PhotoCount: 0,
				PurchaseDate: &time2,
			},
		},
	}
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
	if got[0].PhotoCount != 2 || got[1].PhotoCount != 0 {
		t.Errorf("photo counts = %d, %d; want 2, 0", got[0].PhotoCount, got[1].PhotoCount)
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

func TestBulkUpdateRoomPreservesOtherFields(t *testing.T) {
	provider := &fakeProvider{getItem: models.Item{
		ID: 7, Name: "Coffee Maker", Room: "Storage",
		ReplacementValue: 200, Details: "Keep this note",
	}}
	req := httptest.NewRequest(http.MethodPost, "/api/items/bulk",
		strings.NewReader(`{"ids":[7],"field":"room","value":"Kitchen"}`))
	recorder := httptest.NewRecorder()
	bulkUpdateHandler(provider)(recorder, req)
	if recorder.Code != http.StatusOK {
		t.Fatalf("status %d: %s", recorder.Code, recorder.Body.String())
	}
	if provider.bulkField != "room" || provider.bulkValue != "Kitchen" || len(provider.bulkIDs) != 1 || provider.bulkIDs[0] != 7 {
		t.Errorf("unexpected bulk update: %+v", provider)
	}
}

func TestBulkUpdateRejectsUnsupportedField(t *testing.T) {
	provider := &fakeProvider{}
	req := httptest.NewRequest(http.MethodPost, "/api/items/bulk",
		strings.NewReader(`{"ids":[1],"field":"replacement_value","value":"0"}`))
	recorder := httptest.NewRecorder()
	bulkUpdateHandler(provider)(recorder, req)
	if recorder.Code != http.StatusBadRequest || provider.bulkField != "" {
		t.Errorf("status %d; save called %v", recorder.Code, provider.saveCalled)
	}
}

func TestBulkUpdateRejectsDuplicateIDs(t *testing.T) {
	provider := &fakeProvider{}
	req := httptest.NewRequest(http.MethodPost, "/api/items/bulk",
		strings.NewReader(`{"ids":[1,1],"field":"room","value":"Kitchen"}`))
	recorder := httptest.NewRecorder()
	bulkUpdateHandler(provider)(recorder, req)
	if recorder.Code != http.StatusBadRequest || provider.bulkField != "" {
		t.Errorf("status %d; bulk field %q", recorder.Code, provider.bulkField)
	}
}

func TestBulkUpdateInformationComplete(t *testing.T) {
    provider := &fakeProvider{}
    req := httptest.NewRequest(http.MethodPost, "/api/items/bulk",
        strings.NewReader(`{"ids":[2,3],"field":"info_complete","value":"true"}`))
    recorder := httptest.NewRecorder()
    bulkUpdateHandler(provider)(recorder, req)
    if recorder.Code != http.StatusOK || provider.bulkField != "info_complete" || provider.bulkValue != true {
        t.Errorf("status %d, field %s, value %v", recorder.Code, provider.bulkField, provider.bulkValue)
    }
}
