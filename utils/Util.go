package utils

func InArray(needle int64, haystack []int64) bool {
	for _, i := range haystack {
		if i == needle {
			return true
		}
	}
	return false
}
