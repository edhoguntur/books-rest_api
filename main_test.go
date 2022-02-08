package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllBooks(t *testing.T) {
	req, err := http.NewRequest("GET", "/books", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getBooks)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `[{"id":1,"isbn":"73645","title":"Mudah Kaya","author":{"firstName":"John","lastName":"Doe"}},{"id":2,"isbn":"736343","title":"Mudah Miskin","author":{"firstName":"Jare","lastName":"MbahMu"}},{"id":4,"isbn":"7334334","title":"KonosubaS","author":{"firstName":"CrystaSl","lastName":"MeguminS"}}]
`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
