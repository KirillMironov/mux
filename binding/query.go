package binding

import "net/http"

type QueryParser struct{}

func (QueryParser) TagName() string {
	return "query"
}

func (QueryParser) Lookup(request *http.Request, key string) (string, bool) {
	if !request.URL.Query().Has(key) {
		return "", false
	}
	return request.URL.Query().Get(key), true
}
