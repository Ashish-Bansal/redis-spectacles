package iterator

import "errors"

// ErrEndOfContainer represents that iterator has already reached end of the container
var ErrEndOfContainer = errors.New("Iterator already pointing to end of container")

// Iterator represents anything which can be iterated like other languages
type Iterator interface {
	HasNext() bool
	Next() (interface{}, error)
}

// Iterable interface represents that we can get iterator from this container.
type Iterable interface {
	GetIterator() *Iterator
}

// NewIterator returns Iterator in case it's possible to create one from sent object, otherwise returns error.
func NewIterator(item interface{}) (Iterator, error) {
	switch item.(type) {
	case string:
		return getIterator(item.(string)), nil
	default:
		return nil, errors.New("Don't know how to iterate")
	}
}
