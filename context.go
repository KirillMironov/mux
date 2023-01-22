package mux

import (
	"encoding/json"
	"io"
	"net/http"
)

const (
	mimePlain = "text/plain; charset=utf-8"
	mimeJSON  = "application/json"
)

type Context struct {
	Request *http.Request
	Writer  http.ResponseWriter
	binder  Binder
}

func (c *Context) Status(code int) {
	c.Writer.WriteHeader(code)
}

func (c *Context) String(code int, value string) {
	c.Writer.Header().Set("Content-Type", mimePlain)
	c.Writer.WriteHeader(code)
	_, _ = io.WriteString(c.Writer, value)
}

func (c *Context) JSON(code int, value any) {
	c.Writer.Header().Set("Content-Type", mimeJSON)
	c.Writer.WriteHeader(code)
	_ = json.NewEncoder(c.Writer).Encode(value)
}

func (c *Context) GetQuery(key string) (string, bool) {
	if !c.Request.URL.Query().Has(key) {
		return "", false
	}
	return c.Request.URL.Query().Get(key), true
}

func (c *Context) GetHeader(key string) (string, bool) {
	if _, ok := c.Request.Header[http.CanonicalHeaderKey(key)]; !ok {
		return "", false
	}
	return c.Request.Header.Get(key), true
}

func (c *Context) Bind(target any) error {
	return c.binder.Bind(c.Request, target)
}
