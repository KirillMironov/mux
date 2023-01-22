package mux

import (
	"net/http"
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
