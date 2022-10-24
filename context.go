package beaver

import "net/http"

type Context struct {
	Request *http.Request
	Writer  http.ResponseWriter
}

func (c *Context) String(code int, value string) {
	c.Writer.WriteHeader(code)
	_, _ = c.Writer.Write([]byte(value))
}
