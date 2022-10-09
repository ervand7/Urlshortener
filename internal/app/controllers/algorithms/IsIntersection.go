package algorithms

func IsIntersection(first, second []string) bool {
	hash := make(map[string]bool)
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
