package beaver

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMux(t *testing.T) {
	mux := NewMux()

	mux.Get("/foo", foo)

	mux.Get("/foo/bar", bar)

	checkBody(t, mux, http.MethodGet, "/foo", "foo", http.StatusOK)

	checkBody(t, mux, http.MethodGet, "/foo/bar", "", http.StatusInternalServerError)
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

	checkBody(t, mux, http.MethodGet, "/api/foo", "foo", http.StatusOK)

	checkBody(t, mux, http.MethodGet, "/api/foo/bar", "", http.StatusInternalServerError)

	checkBody(t, mux, http.MethodGet, "/api/v1/foo", "foo", http.StatusOK)

	checkBody(t, mux, http.MethodGet, "/api/v1/foo/bar", "", http.StatusInternalServerError)
}

func TestErrorHandler(t *testing.T) {
	mux := NewMux()

	mux.SetErrorHandler(func(_ error, w http.ResponseWriter) {
		w.WriteHeader(http.StatusNoContent)
		_, _ = w.Write([]byte("error"))
	})

	mux.Get("/foo/bar", bar)

	checkBody(t, mux, http.MethodGet, "/foo/bar", "error", http.StatusNoContent)
}

func foo(c *Context) error {
	c.String(http.StatusOK, "foo")
	return nil
}

func bar(*Context) error {
	return errors.New("bar")
}

func checkBody(t *testing.T, mux *Mux, method, path, expectedBody string, expectedCode int) {
	t.Helper()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, nil)

	mux.ServeHTTP(rec, req)

	if body := rec.Body.String(); body != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, body)
	}

	if rec.Code != expectedCode {
		t.Errorf("expected status code %d, got %d", expectedCode, rec.Code)
	}
}
