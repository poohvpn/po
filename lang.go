package pooh

import (
	"github.com/hashicorp/go-multierror"
	"reflect"
)

type Map map[string]interface{}

func (m *Map) Set(key string, value interface{}) {
	if m == nil {
		return
	}
	if *m == nil {
		*m = make(Map)
	}
	(*m)[key] = value
	return
}

func (m Map) GetBool(key string) (value bool) {
	v, ok := m[key]
	if !ok {
		return
	}
	value, _ = v.(bool)
	return
}

func (m Map) GetInt(key string) (value int) {
	v, ok := m[key]
	if !ok {
		return
	}
	value, _ = v.(int)
	return
}

func (m Map) GetString(key string) (value string) {
	v, ok := m[key]
	if !ok {
		return
	}
	value, _ = v.(string)
	return
}

func Errors(i ...interface{}) (err error) {
	for _, v := range i {
		if e, ok := v.(error); ok && e != nil {
			err = multierror.Append(err, e)
		}
	}
	return
}

func IsNil(object interface{}) bool {
	if object == nil {
		return true
	}

	value := reflect.ValueOf(object)
	switch value.Kind() {
	case
		reflect.Interface,
		reflect.Ptr,
		reflect.Slice,
		reflect.Map,
		reflect.Func,
		reflect.Chan:
		return value.IsNil()
	default:
		return false
	}
}

func Duplicate(bs []byte) []byte {
	if bs == nil {
		return nil
	}
	bsDup := make([]byte, len(bs))
	if len(bs) > 0 {
		copy(bsDup, bs)
	}
	return bsDup
}
