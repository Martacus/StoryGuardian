package utility

func GetAllValues[K comparable, T comparable](m map[K]T) []T {
	values := make([]T, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

func RemoveStringFromSlice(slice []string, str string) []string {
	var result []string
	for _, v := range slice {
		if v != str {
			result = append(result, v)
		}
	}
	return result
}
