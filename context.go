package beaver

import "net/http"

type Context struct {
	request *http.Request
	writer  http.ResponseWriter
}

func (c *Context) String(code int, s string) {
	c.writer.WriteHeader(code)
	_, _ = c.writer.Write([]byte(s))
}

func (c *Context) Request() *http.Request {
	return c.request
}

func (c *Context) Writer() http.ResponseWriter {
	return c.writer
}
