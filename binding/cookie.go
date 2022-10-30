package binding

import "net/http"

type CookieParser struct{}

func (CookieParser) TagName() string {
	return "cookie"
}

func (CookieParser) Lookup(request *http.Request, key string) (string, bool) {
	cookie, err := request.Cookie(key)
	if err != nil {
		return "", false
	}
	return cookie.Value, true
}
