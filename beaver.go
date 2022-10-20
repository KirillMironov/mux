package beaver

import "net/http"

type (
	HandlerFunc func(*Context) error

	ErrorHandler func(error, http.ResponseWriter)

	route struct {
		method      string
		path        string
		handlerFunc HandlerFunc
	}
)

var DefaultErrorHandler = func(_ error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}
