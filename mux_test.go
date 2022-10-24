package beaver

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMux(t *testing.T) {
	mux := NewMux()

	mux.Get("/foo", foo)

	mux.Get("/foo/bar", bar)

	makeRequest(t, "/foo", http.MethodGet, "foo", http.StatusOK, mux)

	makeRequest(t, "/foo/bar", http.MethodGet, "", http.StatusInternalServerError, mux)
}

func TestGroup(t *testing.T) {
	mux := NewMux()

	api := mux.Group("/api")
	{
		api.Get("/foo", foo)
		api.Get("/foo/bar", bar)

		v1 := api.Group("/v1")
		{
			v1.Get("/foo", foo)
			v1.Get("/foo/bar", bar)
		}
	}

	makeRequest(t, "/api/foo", http.MethodGet, "foo", http.StatusOK, mux)

	makeRequest(t, "/api/foo/bar", http.MethodGet, "", http.StatusInternalServerError, mux)

	makeRequest(t, "/api/v1/foo", http.MethodGet, "foo", http.StatusOK, mux)

	makeRequest(t, "/api/v1/foo/bar", http.MethodGet, "", http.StatusInternalServerError, mux)
}

func foo(c *Context) error {
	c.String(http.StatusOK, "foo")
	return nil
}

func bar(*Context) error {
	return errors.New("bar")
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
