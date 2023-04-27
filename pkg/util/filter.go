package util

func Filter[T comparable](elems []T, match func(e T) bool) []T {
	var out []T
	for _, e := range elems {
		if match(e) {
			out = append(out, e)
		}
	}
	return out
}
