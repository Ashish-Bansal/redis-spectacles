package addable

import (
	"errors"
	"reflect"
)

// Addable is the interface which when implemented by structs, they can be added together into single object
// For basic types, add logic inside below given Add func.
type Addable interface {
	Add(Addable) (Addable, error)
}

// Add returns Iterator in case it's possible to create one from sent object, otherwise returns error.
func Add(a interface{}, b interface{}) (interface{}, error) {
	if a == nil {
		return b, nil
	}

	if b == nil {
		return a, nil
	}

	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return nil, errors.New("Add fn called on interfaces of two different types")
	}

	aAddable, ok := a.(Addable)
	bAddable, ok := b.(Addable)
	if ok {
		return aAddable.Add(bAddable)
	}

	switch a.(type) {
	case string:
		return a.(string) + b.(string), nil
	default:
		return nil, errors.New("Don't know how to iterate")
	}
}
