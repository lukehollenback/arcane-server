package util

//
// SliceContainsString searches for the specified needle in the specified haystack.
//
func SliceContainsString(needle string, haystack []string) bool {
	for _, e := range haystack {
		if e == needle {
			return true
		}
	}

	return false
}

//
// SliceContainsInt searches for the specified needle in the specified haystack.
//
func SliceContainsInt(needle int, haystack []int) bool {
	for _, e := range haystack {
		if e == needle {
			return true
		}
	}

	return false
}
