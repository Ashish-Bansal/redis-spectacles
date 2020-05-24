package utils

// Min - a super complex function which golang authors never thought of implementing.
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
