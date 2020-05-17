package comparable

import "testing"

func TestStringCompare(t *testing.T) {
	input := [][]string{
		{"a", "b"},
		{"b", "a"},
		{"", "b"},
		{"ab", "abc"},
	}
	output := []bool{
		true,
		false,
		true,
		true,
	}

	for index := range input {
		testcase := input[index]
		expectedOutput := output[index]

		firstString := testcase[0]
		secondString := testcase[1]
		result, err := LessThan(firstString, secondString)
		if err != nil {
			t.Errorf("%v", err)
		}

		if expectedOutput != result {
			t.Errorf(
				"LessThan(%s, %s) - Expected %v, got %v",
				firstString,
				secondString,
				expectedOutput,
				result,
			)
		}
	}
}
