package pooh

import (
	"net"
	"reflect"
)

func IsIPv4(ip net.IP) bool {
	return ip.To4() != nil
}

func IsIPv6(ip net.IP) bool {
	return ip.To16() != nil && !IsIPv4(ip)
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
