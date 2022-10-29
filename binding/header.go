package binding

import "net/http"

type HeaderParser struct{}

func (HeaderParser) TagName() string {
	return "header"
}

func (HeaderParser) Lookup(request *http.Request, key string) (string, bool) {
	if _, ok := request.Header[http.CanonicalHeaderKey(key)]; !ok {
		return "", false
	}
	return request.Header.Get(key), true
}
