package iterator

type stringIterator struct {
	str   []rune
	index int
}

func (it *stringIterator) HasNext() bool {
	return it.index < len(it.str)
}

func (it *stringIterator) Next() (interface{}, error) {
	if !it.HasNext() {
		return nil, ErrEndOfContainer
	}
	it.index++
	return it.str[it.index-1], nil
}

func getIterator(str string) Iterator {
	return &stringIterator{str: []rune(str)}
}
