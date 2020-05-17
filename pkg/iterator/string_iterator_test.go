package iterator

import (
	"fmt"
	"testing"
)

func ExampleNewIterator() {
	testString := "test"
	it := getIterator(testString)
	for it.HasNext() {
		character, _ := it.Next()
		fmt.Printf("%s", character)
	}
	// Output:
	// test
}

func TestIteratorValues(t *testing.T) {
	testString := "test"
	testStringChars := []string{"t", "e", "s", "t"}

	testStringCharsIndex := 0
	it := getIterator(testString)
	for it.HasNext() {
		character, _ := it.Next()
		if character != testStringChars[testStringCharsIndex] {
			t.Errorf(
				"Iterator values didn't match. Expected %s, got %s.",
				testStringChars[testStringCharsIndex],
				character,
			)
		}
		testStringCharsIndex++
	}
}
