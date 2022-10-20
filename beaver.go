package beaver

import "net/http"

type (
	HandlerFunc func(http.ResponseWriter, *http.Request) error

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
