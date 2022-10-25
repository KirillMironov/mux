package beaver

import (
	"errors"
	"net/http"
	"testing"
)

func TestMux(t *testing.T) {
	mux := NewMux()

	mux.Get("/foo", foo)

	mux.Get("/foo/bar", bar)

	testCase{
		mux:    mux,
		method: http.MethodGet,
		path:   "/foo",
		expected: expected{
			statusCode:  http.StatusOK,
			contentType: mimePlain,
			body:        "foo",
		},
	}.run(t)

	testCase{
		mux:    mux,
		method: http.MethodGet,
		path:   "/foo/bar",
		expected: expected{
			statusCode:  http.StatusInternalServerError,
			contentType: "",
			body:        "",
		},
	}.run(t)
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

	testCase{
		mux:    mux,
		method: http.MethodGet,
		path:   "/api/foo",
		expected: expected{
			statusCode:  http.StatusOK,
			contentType: mimePlain,
			body:        "foo",
		},
	}.run(t)

	testCase{
		mux:    mux,
		method: http.MethodGet,
		path:   "/api/foo/bar",
		expected: expected{
			statusCode:  http.StatusInternalServerError,
			contentType: "",
			body:        "",
		},
	}.run(t)

	testCase{
		mux:    mux,
		method: http.MethodGet,
		path:   "/api/v1/foo",
		expected: expected{
			statusCode:  http.StatusOK,
			contentType: mimePlain,
			body:        "foo",
		},
	}.run(t)

	testCase{
		mux:    mux,
		method: http.MethodGet,
		path:   "/api/v1/foo/bar",
		expected: expected{
			statusCode:  http.StatusInternalServerError,
			contentType: "",
			body:        "",
		},
	}.run(t)
}

func TestErrorHandler(t *testing.T) {
	mux := NewMux()

	mux.SetErrorHandler(func(_ error, w http.ResponseWriter) {
		w.WriteHeader(http.StatusNoContent)
		_, _ = w.Write([]byte("error"))
	})

	mux.Get("/foo/bar", bar)

	testCase{
		mux:    mux,
		method: http.MethodGet,
		path:   "/foo/bar",
		expected: expected{
			statusCode:  http.StatusNoContent,
			contentType: "",
			body:        "error",
		},
	}.run(t)
}

func foo(c *Context) error {
	c.String(http.StatusOK, "foo")
	return nil
}

func bar(*Context) error {
	return errors.New("bar")
}
