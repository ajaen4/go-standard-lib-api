package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ajaen4/go-standard-lib-api/pkg/api_errors"
)

func TestHealthCheck(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	HealthCheck(w, req)

	expectedCode := http.StatusOK
	if status := w.Code; status != expectedCode {
		t.Errorf(
			"handler returned wrong status code: got %v want %v",
			status,
			expectedCode,
		)
	}

	expectedCType := "text/plain; charset=utf-8"
	if contentType := w.Header().Get("Content-Type"); contentType != expectedCType {
		t.Errorf(
			"handler returned wrong content type: got %s want %s",
			contentType,
			expectedCType,
		)
	}

	expectedResp := "OK"
	if response := w.Body.String(); response != expectedResp {
		t.Errorf(
			"handler returned wrong plain test code: got %s want %s",
			response,
			expectedResp,
		)
	}
}

func TestNewHandler(t *testing.T) {
	OKHandler := func(w http.ResponseWriter, r *http.Request) error {
		return nil
	}

	OKreq, OKerr := http.NewRequest("GET", "/api/test", nil)
	if OKerr != nil {
		t.Fatal(OKerr)
	}

	Okw := httptest.NewRecorder()
	NewHandler(OKHandler).ServeHTTP(Okw, OKreq)
	if Okw.Body.Len() != 0 {
		t.Errorf(
			"handler returned wrong error response: got %v want an empty response",
			Okw.Body,
		)
	}

	KOHandler := func(w http.ResponseWriter, r *http.Request) error {
		return errors.New("error message")
	}

	KOreq, KOerr := http.NewRequest("GET", "/api/test", nil)
	if KOerr != nil {
		t.Fatal(KOerr)
	}

	w := httptest.NewRecorder()
	KOExpectedResp := api_errors.InternalErr{
		HttpCode: http.StatusInternalServerError,
		Message:  "internal server error",
	}
	KOresponse := api_errors.InternalErr{}
	NewHandler(KOHandler).ServeHTTP(w, KOreq)
	if err := json.NewDecoder(w.Body).Decode(&KOresponse); err != nil {
		t.Fatal(err)
	}

	if KOresponse != KOExpectedResp {
		t.Errorf(
			"handler returned wrong error response: got %v want %v",
			KOresponse,
			KOExpectedResp,
		)
	}

}
