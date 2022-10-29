package binding

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
)

type Binder struct {
	parsers []Parser
}

type Parser interface {
	TagName() string
	Lookup(request *http.Request, key string) (string, bool)
}

func NewBinder() *Binder {
	return &Binder{
		parsers: []Parser{
			&CookieParser{},
			&HeaderParser{},
			&QueryParser{},
		},
	}
}

func (b *Binder) Bind(request *http.Request, target any) (err error) {
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

	err = bindBody(request, target)
	if err != nil {
		return err
	}

	for i := 0; i < targetType.NumField(); i++ {
		tag := targetType.Field(i).Tag

		for _, parser := range b.parsers {
			tagValue, ok := tag.Lookup(parser.TagName())
			if !ok || tagValue == "" {
				continue
			}

			value, ok := parser.Lookup(request, tagValue)
			if !ok {
				continue
			}

			err = setField(value, targetValue.Field(i))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func setField(value string, targetField reflect.Value) error {
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
