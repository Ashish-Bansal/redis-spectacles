package addable

import "testing"

func TestStringAddition(t *testing.T) {
	testSuite := [][]string{{"a", "b", "ab"}, {"", "a", "a"}}
	for _, testcase := range testSuite {
		firstString := testcase[0]
		secondString := testcase[1]
		expectedString := testcase[2]
		result, err := Add(firstString, secondString)
		if err != nil {
			t.Errorf("%v", err)
		}

		if expectedString != result {
			t.Errorf(
				"Add(%s, %s) - Expected %s, got %s",
				firstString,
				secondString,
				expectedString,
				result,
			)
		}
	}
}
