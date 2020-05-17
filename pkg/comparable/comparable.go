package comparable

import (
	"errors"
	"reflect"
)

// Comparable is the interface which when implemented by structs, can be
// compared with each other.
// For basic types, add logic inside below given LessThan/Equal func.
type Comparable interface {
	LessThan(Comparable) (bool, error)
	Equal(Comparable) (bool, error)
}

// LessThan returns whether a is less than b or not.
// Returns error in case it's not possible to compare two interfaces.
func LessThan(a interface{}, b interface{}) (bool, error) {
	if a == nil {
		return true, nil
	}

	if b == nil {
		return false, nil
	}

	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false, errors.New("LessThan fn called on interfaces of two different types")
	}

	aComparable, ok := a.(Comparable)
	bComparable, ok := b.(Comparable)
	if ok {
		return aComparable.LessThan(bComparable)
	}

	switch a.(type) {
	case string:
		return a.(string) < b.(string), nil
	default:
		return false, errors.New("Don't know how to iterate")
	}
}

// Equal returns whether two elements are equal or not
func Equal(a interface{}, b interface{}) (bool, error) {
	isLess, err := LessThan(a, b)
	if isLess || err != nil {
		return false, err
	}

	isGreater, err := LessThan(b, a)
	if isGreater || err != nil {
		return false, err
	}

	return true, nil
}
