package beaver

import "net/http"

type Mux struct {
	routes       map[string]route
	errorHandler ErrorHandler
}

func NewMux() *Mux {
	return &Mux{
		routes:       make(map[string]route),
		errorHandler: DefaultErrorHandler,
	}
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route, ok := m.routes[r.RequestURI]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if route.method != r.Method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := route.handlerFunc(w, r)
	if err != nil {
		m.errorHandler(err, w)
	}
}

func (m *Mux) Get(path string, hf HandlerFunc) {
	m.routes[path] = route{
		http.MethodGet,
		path,
		hf,
	}
}
