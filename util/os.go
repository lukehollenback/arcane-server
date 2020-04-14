package util

import "os"

//
// GetEnv attempts to get a value from the environment for the specified variable. If the specified
// variable is unset, simply returns the specified default value.
//
func GetEnv(name string, value string) string {
	fndVal, fnd := os.LookupEnv(name)
	if !fnd {
		return value
	}

	return fndVal
}
