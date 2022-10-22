package beaver

import "net/http"

var DefaultErrorHandler = func(_ error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}

type (
	HandlerFunc func(*Context) error

	ErrorHandler func(error, http.ResponseWriter)

	errorHandlerFunc struct {
		handlerFunc  HandlerFunc
		errorHandler ErrorHandler
		method       string
	}
)

func (ehf *errorHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != ehf.method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var context = &Context{
		request: r,
		writer:  w,
	}

	err := ehf.handlerFunc(context)
	if err != nil {
		ehf.errorHandler(err, w)
	}
}
