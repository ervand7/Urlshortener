package algorithms

// Issubset checks whether first contains second
func Issubset(first, second []string) bool {
	hash := make(map[string]bool, len(first))
	for _, i := range first {
		hash[i] = true
	}
	for _, i := range second {
		if _, ok := hash[i]; !ok {
			return false
		}
	}

	return true
}
