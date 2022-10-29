package binding

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
)

const (
	mimeJSON    = "application/json"
	mimeXML     = "application/xml"
	mimeTextXML = "text/xml"
)

func bindBody(request *http.Request, target any) error {
	switch request.Header.Get("Content-Type") {
	case mimeJSON:
		err := json.NewDecoder(request.Body).Decode(target)
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
	case mimeXML, mimeTextXML:
		err := xml.NewDecoder(request.Body).Decode(target)
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
	}

	return nil
}
