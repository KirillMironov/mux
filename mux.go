package beaver

import "net/http"

type Mux struct {
	handler      *http.ServeMux
	errorHandler ErrorHandler
}

func NewMux() *Mux {
	return &Mux{
		handler:      http.DefaultServeMux,
		errorHandler: DefaultErrorHandler,
	}
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.handler.ServeHTTP(w, r)
}

func (m *Mux) Get(path string, hf HandlerFunc) {
	m.wrap(http.MethodGet, path, hf)
}

func (m *Mux) Head(path string, hf HandlerFunc) {
	m.wrap(http.MethodHead, path, hf)
}

func (m *Mux) Post(path string, hf HandlerFunc) {
	m.wrap(http.MethodPost, path, hf)
}

func (m *Mux) Put(path string, hf HandlerFunc) {
	m.wrap(http.MethodPut, path, hf)
}

func (m *Mux) Patch(path string, hf HandlerFunc) {
	m.wrap(http.MethodPatch, path, hf)
}

func (m *Mux) Delete(path string, hf HandlerFunc) {
	m.wrap(http.MethodDelete, path, hf)
}

func (m *Mux) Connect(path string, hf HandlerFunc) {
	m.wrap(http.MethodConnect, path, hf)
}

func (m *Mux) Options(path string, hf HandlerFunc) {
	m.wrap(http.MethodOptions, path, hf)
}

func (m *Mux) Trace(path string, hf HandlerFunc) {
	m.wrap(http.MethodTrace, path, hf)
}

func (m *Mux) wrap(method, path string, hf HandlerFunc) {
	m.handler.Handle(path, errorHandlerFunc{
		handlerFunc:  hf,
		errorHandler: m.errorHandler,
		method:       method,
	})
}
