package internal

import (
	"reflect"
)

func Zeroed[T any] () T {
	return *new(T) ; 
}

func isNilInterface(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return rv.IsNil()
	default:
		return false
	}
}

func NMutCast[T any]( from any ) T {
	if !isNilInterface(from) {
		if from, ok := any(from).(T) ; ok {
			return from ; 
		}
		return Zeroed[T]() ; 
	}
	return Zeroed[T]() ; 
}


func MutCast[T any]( from *any ) {
	if r, ok := any(from).(T) ; ok {
		*from = r ; 
	}
	*from = Zeroed[T]()
}





