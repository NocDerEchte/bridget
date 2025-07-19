package helper

import "os"

// GetEnv returns the value of a given environment variable (key). If the result is empty returns fallback.
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
