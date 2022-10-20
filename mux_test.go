package beaver

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMux(t *testing.T) {
	mux := NewMux()

	mux.Get("/foo", func(w http.ResponseWriter, r *http.Request) error {
		fmt.Fprint(w, "foo")
		return nil
	})

	mux.Get("/foo/bar", func(w http.ResponseWriter, r *http.Request) error {
		fmt.Fprint(w, "bar")
		return nil
	})

	mux.Get("/error", func(w http.ResponseWriter, r *http.Request) error {
		return errors.New("error")
	})

	makeRequest(t, "/foo", http.MethodGet, "foo", http.StatusOK, mux)

	makeRequest(t, "/foo/bar", http.MethodGet, "bar", http.StatusOK, mux)

	makeRequest(t, "/error", http.MethodGet, "", http.StatusInternalServerError, mux)
}

func makeRequest(t *testing.T, path, method, expectedBody string, expectedStatusCode int, mux *Mux) {
	t.Helper()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, nil)

	mux.ServeHTTP(rec, req)

	body, err := io.ReadAll(rec.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != expectedBody {
		t.Fatalf("expected %q, got %q", expectedBody, string(body))
	}

	if rec.Code != expectedStatusCode {
		t.Fatalf("expected status code %d, got %d", expectedStatusCode, rec.Code)
	}
}
