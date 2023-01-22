package mux

import "net/http"

var DefaultErrorHandler = func(_ error, c *Context) {
	c.Status(http.StatusInternalServerError)
}

type (
	HandlerFunc func(*Context) error

	ErrorHandler func(error, *Context)

	Binder interface {
		Bind(request *http.Request, target any) error
	}

	wrapper struct {
		handlerFunc  HandlerFunc
		errorHandler ErrorHandler
		binder       Binder
		method       string
	}
)

func (wrapper *wrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != wrapper.method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var context = &Context{
		Request: r,
		Writer:  w,
		binder:  wrapper.binder,
	}

	err := wrapper.handlerFunc(context)
	if err != nil {
		wrapper.errorHandler(err, context)
	}
}
