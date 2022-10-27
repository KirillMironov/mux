package beaver

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
)

const (
	mimePlain   = "text/plain; charset=utf-8"
	mimeJSON    = "application/json"
	mimeXML     = "application/xml"
	mimeTextXML = "text/xml"
)

const (
	queryTag  = "query"
	headerTag = "header"
)

type Context struct {
	Request *http.Request
	Writer  http.ResponseWriter
}

func (c *Context) Status(code int) {
	c.Writer.WriteHeader(code)
}

func (c *Context) String(code int, value string) {
	c.Writer.Header().Set("Content-Type", mimePlain)
	c.Writer.WriteHeader(code)
	_, _ = io.WriteString(c.Writer, value)
}

func (c *Context) JSON(code int, value any) {
	c.Writer.Header().Set("Content-Type", mimeJSON)
	c.Writer.WriteHeader(code)
	_ = json.NewEncoder(c.Writer).Encode(value)
}

func (c *Context) GetQuery(key string) (string, bool) {
	if !c.Request.URL.Query().Has(key) {
		return "", false
	}
	return c.Request.URL.Query().Get(key), true
}

func (c *Context) GetHeader(key string) (string, bool) {
	if _, ok := c.Request.Header[http.CanonicalHeaderKey(key)]; !ok {
		return "", false
	}
	return c.Request.Header.Get(key), true
}

func (c *Context) Bind(target any) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("unable to bind: %v", r)
		}
	}()

	targetType := reflect.TypeOf(target)

	if targetType.Kind() != reflect.Ptr || targetType.Elem().Kind() != reflect.Struct {
		return errors.New("target must be a pointer to a struct")
	}

	targetType = targetType.Elem()

	targetValue := reflect.ValueOf(target).Elem()

	switch c.Request.Header.Get("Content-Type") {
	case mimeJSON:
		err = json.NewDecoder(c.Request.Body).Decode(target)
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
	case mimeXML, mimeTextXML:
		err = xml.NewDecoder(c.Request.Body).Decode(target)
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
	}

	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i)

		err = c.bind(queryTag, field.Tag, targetValue.Field(i), c.GetQuery)
		if err != nil && !errors.Is(err, errTagNotFound) && !errors.Is(err, errValueNotFound) {
			return err
		}

		err = c.bind(headerTag, field.Tag, targetValue.Field(i), c.GetHeader)
		if err != nil && !errors.Is(err, errTagNotFound) && !errors.Is(err, errValueNotFound) {
			return err
		}
	}

	return nil
}

type getter func(key string) (string, bool)

var (
	errTagNotFound   = errors.New("tag not found")
	errValueNotFound = errors.New("value not found")
)

func (c *Context) bind(tagKey string, tag reflect.StructTag, targetField reflect.Value, get getter) error {
	tagValue, ok := tag.Lookup(tagKey)
	if !ok || tagValue == "" {
		return errTagNotFound
	}

	value, ok := get(tagValue)
	if !ok {
		return errValueNotFound
	}

	switch targetField.Type().Kind() {
	case reflect.String:
		targetField.SetString(value)
	case reflect.Bool:
		v, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		targetField.SetBool(v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		targetField.SetInt(v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		targetField.SetUint(v)
	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		targetField.SetFloat(v)
	default:
		return fmt.Errorf("unsupported type: %v", targetField.Type().Kind())
	}

	return nil
}
