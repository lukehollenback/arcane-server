package util

//
// GetStrVal returns either the provided value, or the specified default value if the provided
// value is empty.
//
func GetStrVal(val string, def string) string {
	if len(val) == 0 {
		return def
	}

	return val
}
