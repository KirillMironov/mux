package beaver

import (
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

	makeRequest(t, "/foo", http.MethodGet, "foo", mux)

	makeRequest(t, "/foo/bar", http.MethodGet, "bar", mux)
}

func makeRequest(t *testing.T, path, method, expectedBody string, mux *Mux) {
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
}
