package util

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

const (
	maxResourceNameLength    = 253
	resourceNameRegexpString = "([A-Za-z0-9][-A-Za-z0-9_.]*)?"
)

var (
	resourceNameRegexp = regexp.MustCompile(resourceNameRegexpString)
)

func ValidateResourceName(name string) error {
	if len(name) > maxResourceNameLength {
		return fmt.Errorf("invalid resource name: length %d is bigger than maximum %d", len(name), maxResourceNameLength)
	}

	if !resourceNameRegexp.MatchString(name) {
		return fmt.Errorf("invalid resource name: it must match regexp '%s'", resourceNameRegexpString)
	}

	return nil
}

type ExcessKeysError struct {
	Path       string
	ExcessKeys []string
}

func (e ExcessKeysError) Error() string {
	return fmt.Sprintf("excess keys in %q: %#v", e.Path, e.ExcessKeys)
}

func NewExcessKeysErrorFromMap(path string, m map[string]interface{}) ExcessKeysError {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	return ExcessKeysError{Path: path, ExcessKeys: keys}
}

func IsRequiredField(f reflect.StructField) bool {
	tag := f.Tag.Get("yaml")
	options := strings.Split(tag, ",")
	for _, o := range options {
		if o == "omitempty" {
			return false
		}
	}
	return true
}

func IsUnset(v reflect.Value) bool {
	if v.Kind() == reflect.Slice && v.Len() == 0 {
		return true
	}

	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}

func isTraversable(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Struct:
		return true
	case reflect.Slice, reflect.Array:
		return isTraversable(t.Elem())
	default:
		return false
	}
}

func ValidateRequiredFields(i interface{}) error {
	var v reflect.Value
	var ok bool
	if v, ok = i.(reflect.Value); !ok {
		v = reflect.ValueOf(i)
	}

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch t := v.Kind(); t {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fv := v.Field(i)
			ft := v.Type().Field(i)

			if IsRequiredField(ft) {
				// zero value means it wasn't set
				// in cases when the zero value is valid value you need to make it a pointer to that value
				if IsUnset(fv) {
					return fmt.Errorf("Required field %q has unset value %#v", ft.Name, fv)
				}
			}

			if isTraversable(fv.Type()) {
				err := ValidateRequiredFields(fv)
				if err != nil {
					return err
				}
			}
		}
	case reflect.Slice, reflect.Array:
		if isTraversable(v.Type()) {
			for i := 0; i < v.Len(); i++ {
				err := ValidateRequiredFields(v.Index(i))
				if err != nil {
					return err
				}
			}
		}
	default:
		return fmt.Errorf("unsupported type %q", t)
	}

	return nil
}
