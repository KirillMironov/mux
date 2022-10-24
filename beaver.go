package beaver

import "net/http"

var DefaultErrorHandler = func(_ error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}

type (
	HandlerFunc func(*Context) error

	ErrorHandler func(error, http.ResponseWriter)

	wrapper struct {
		handlerFunc  HandlerFunc
		errorHandler ErrorHandler
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
	}

	err := wrapper.handlerFunc(context)
	if err != nil {
		wrapper.errorHandler(err, w)
	}
}
