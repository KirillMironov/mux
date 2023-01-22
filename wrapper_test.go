package mux

import (
	"net/http/httptest"
	"testing"
)

type (
	testCase struct {
		mux    *Mux
		method string
		path   string
		expected
	}

	expected struct {
		statusCode  int
		contentType string
		body        string
	}
)

func (tc testCase) run(t *testing.T) {
	t.Helper()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(tc.method, tc.path, nil)

	tc.mux.ServeHTTP(rec, req)

	if rec.Code != tc.expected.statusCode {
		t.Errorf("expected status code %d, got %d", tc.expected.statusCode, rec.Code)
	}

	if contentType := rec.Header().Get("Content-Type"); contentType != tc.expected.contentType {
		t.Errorf("expected Content-Type %q, got %q", tc.expected.contentType, contentType)
	}

	if body := rec.Body.String(); body != tc.expected.body {
		t.Errorf("expected body %q, got %q", tc.expected.body, body)
	}
}
