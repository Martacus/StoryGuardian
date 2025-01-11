package utility

func GetAllValues[K comparable, T comparable](m map[K]T) []T {
	values := make([]T, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}
