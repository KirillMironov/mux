package beaver

import (
	"net/http"
	"net/http/httptest"
	"strings"
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

	testCase{
		mux:    mux,
		method: http.MethodGet,
		path:   "/foo",
		expected: expected{
			statusCode:  http.StatusUnauthorized,
			contentType: "",
			body:        "",
		},
	}.run(t)

	testCase{
		mux:    mux,
		method: http.MethodGet,
		path:   "/bar",
		expected: expected{
			statusCode:  http.StatusOK,
			contentType: "",
			body:        "",
		},
	}.run(t)
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

	testCase{
		mux:    mux,
		method: http.MethodGet,
		path:   "/foo",
		expected: expected{
			statusCode:  http.StatusCreated,
			contentType: mimePlain,
			body:        "foo",
		},
	}.run(t)

	testCase{
		mux:    mux,
		method: http.MethodGet,
		path:   "/bar",
		expected: expected{
			statusCode:  http.StatusNoContent,
			contentType: mimePlain,
			body:        "",
		},
	}.run(t)
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

	testCase{
		mux:    mux,
		method: http.MethodGet,
		path:   "/foo",
		expected: expected{
			statusCode:  http.StatusCreated,
			contentType: mimeJSON,
			body:        "{\"foo\":\"bar\"}\n",
		},
	}.run(t)

	testCase{
		mux:    mux,
		method: http.MethodGet,
		path:   "/bar",
		expected: expected{
			statusCode:  http.StatusNoContent,
			contentType: mimeJSON,
			body:        "null\n",
		},
	}.run(t)
}

func TestContext_Bind(t *testing.T) {
	type form struct {
		Username    string `query:"username"`
		Id          uint   `query:"id"`
		AccessToken string `header:"access-token"`
		Email       string `json:"email"`
	}

	req := httptest.NewRequest(http.MethodGet, "/?username=admin&id=55", strings.NewReader(`{"email": "me"}`))
	req.Header.Set("Content-Type", mimeJSON)
	req.Header.Set("access-token", "8888")

	var (
		target form

		expectedResult = form{
			Username:    "admin",
			Id:          55,
			AccessToken: "8888",
			Email:       "me",
		}

		context = &Context{Request: req}
	)

	err := context.Bind(&target)
	if err != nil {
		t.Error(err)
	}

	if target != expectedResult {
		t.Errorf("expected %+v, got %+v", expectedResult, target)
	}
}
