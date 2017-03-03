package collectionutils

func MapFromArray(array []string) map[string]bool {
	m := map[string]bool{}
	for _, element := range array {
		m[element] = true
	}
	return m
}
