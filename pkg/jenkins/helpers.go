package jenkins

// InIntSlice ...
func InIntSlice(slice []int, elem int) bool {
	found := false
	for _, i := range slice {
		if i == elem {
			found = true
			break
		}
	}
	return found
}

// RemoveIndexes ...
func RemoveIndexes(slice []string, indexes []int) []string {
	newSlice := []string{}
	for i, s := range slice {
		if !InIntSlice(indexes, i) {
			newSlice = append(newSlice, s)
		}
	}
	return newSlice
}
