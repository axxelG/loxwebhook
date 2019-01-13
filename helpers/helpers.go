package helpers

// IsStringInSlice returns true if str is in list
func IsStringInSlice(str string, list []string) bool {
	for _, i := range list {
		if i == str {
			return true
		}
	}
	return false
}

// GetMapStringKeyFromStringValue returns the key of a map[string]string item with value str
func GetMapStringKeyFromStringValue(str string, m map[string]string) (key string, ok bool) {
	for k, v := range m {
		if v == str {
			key = k
			ok = true
			break
		}
	}
	return
}
