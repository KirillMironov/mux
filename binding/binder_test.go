package binding

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBind(t *testing.T) {
	type form struct {
		Id          uint   `query:"id"`
		Name        string `query:"name"`
		Email       string `json:"email"`
		AccessToken int    `header:"Authorization"`
	}

	req := httptest.NewRequest(http.MethodGet, "/?id=55&name=admin", strings.NewReader(`{"email": "x@x.com"}`))
	req.Header.Set("Content-Type", mimeJSON)
	req.Header.Set("Authorization", "8888")

	var (
		binder = NewBinder()

		target form

		expected = form{
			Id:          55,
			Name:        "admin",
			Email:       "x@x.com",
			AccessToken: 8888,
		}
	)

	err := binder.Bind(req, &target)
	if err != nil {
		t.Error(err)
	}

	if target != expected {
		t.Errorf("expected: %v, got: %v", expected, target)
	}
}
