package beaver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestContext_Status(t *testing.T) {
	mux := NewMux()

	mux.Get("/foo", func(c *Context) error {
		c.Status(http.StatusUnauthorized)
		return nil
	})

	mux.Get("/bar", func(c *Context) error {
		return nil
	})

	checkHeaders(t, mux, http.MethodGet, "/foo", "", "", http.StatusUnauthorized)

	checkHeaders(t, mux, http.MethodGet, "/bar", "", "", http.StatusOK)
}

func TestContext_String(t *testing.T) {
	mux := NewMux()

	mux.Get("/foo", func(c *Context) error {
		c.String(http.StatusCreated, "foo")
		return nil
	})

	mux.Get("/bar", func(c *Context) error {
		c.String(http.StatusNoContent, "")
		return nil
	})

	checkHeaders(t, mux, http.MethodGet, "/foo", "foo", mimePlain, http.StatusCreated)

	checkHeaders(t, mux, http.MethodGet, "/bar", "", mimePlain, http.StatusNoContent)
}

func TestContext_JSON(t *testing.T) {
	mux := NewMux()

	mux.Get("/foo", func(c *Context) error {
		c.JSON(http.StatusCreated, map[string]string{"foo": "bar"})
		return nil
	})

	mux.Get("/bar", func(c *Context) error {
		c.JSON(http.StatusNoContent, nil)
		return nil
	})

	checkHeaders(t, mux, http.MethodGet, "/foo", "{\"foo\":\"bar\"}\n", mimeJSON, http.StatusCreated)

	checkHeaders(t, mux, http.MethodGet, "/bar", "null\n", mimeJSON, http.StatusNoContent)
}

func checkHeaders(t *testing.T, mux *Mux, method, path, expectedBody, expectedContentType string, expectedCode int) {
	t.Helper()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, nil)

	mux.ServeHTTP(rec, req)

	if body := rec.Body.String(); body != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, body)
	}

	if contentType := rec.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("expected Content-Type %q, got %q", expectedContentType, contentType)
	}

	if rec.Code != expectedCode {
		t.Errorf("expected status code %d, got %d", expectedCode, rec.Code)
	}
}
