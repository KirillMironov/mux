package beaver

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
}

func (c *Context) Status(code int) {
	c.Writer.WriteHeader(code)
}

func (c *Context) String(code int, value string) {
	c.setHeaders(code, mimePlain)
	_, _ = io.WriteString(c.Writer, value)
}

func (c *Context) JSON(code int, value any) {
	c.setHeaders(code, mimeJSON)
	_ = json.NewEncoder(c.Writer).Encode(value)
}

func (c *Context) setHeaders(code int, contentType string) {
	c.Writer.Header().Set("Content-Type", contentType)
	c.Writer.WriteHeader(code)
}
