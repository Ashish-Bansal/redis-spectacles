package iterator

import (
	"fmt"
)

func ExampleNewIterator() {
	testString := "test"
	it := getIterator(testString)
	for it.HasNext() {
		character, _ := it.Next()
		fmt.Printf("%c", character)
	}
	// Output:
	// test
}
