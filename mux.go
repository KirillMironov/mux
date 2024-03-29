package mux

import (
	"net/http"

	"github.com/KirillMironov/mux/binding"
)

type Mux struct {
	handler      *http.ServeMux
	errorHandler ErrorHandler
	binder       Binder
	basePath     string
}

func NewMux() *Mux {
	return &Mux{
		handler:      new(http.ServeMux),
		errorHandler: DefaultErrorHandler,
		binder:       binding.NewBinder(),
	}
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.handler.ServeHTTP(w, r)
}

func (m *Mux) SetErrorHandler(eh ErrorHandler) {
	m.errorHandler = eh
}

func (m *Mux) SetBinder(b Binder) {
	m.binder = b
}

func (m *Mux) Group(basePath string) *Mux {
	return &Mux{
		handler:      m.handler,
		errorHandler: m.errorHandler,
		binder:       m.binder,
		basePath:     m.basePath + basePath,
	}
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
	m.handler.Handle(m.basePath+path, &wrapper{
		handlerFunc:  hf,
		errorHandler: m.errorHandler,
		binder:       m.binder,
		method:       method,
	})
}
