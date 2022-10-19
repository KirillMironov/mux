package beaver

import "net/http"

type (
	HandlerFunc func(http.ResponseWriter, *http.Request) error

	ErrorHandler func(error)

	route struct {
		method      string
		path        string
		handlerFunc HandlerFunc
	}
)
