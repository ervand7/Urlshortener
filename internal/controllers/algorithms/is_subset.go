package algorithms

// Issubset checks whether first contains second
func Issubset(first, second []string) bool {
	if len(second) > len(first) {
		return false
	}

	hash := make(map[string]struct{}, len(first))
	for _, i := range first {
		hash[i] = struct{}{}
	}
	for _, i := range second {
		if _, ok := hash[i]; !ok {
			return false
		}
	}

	return true
}
