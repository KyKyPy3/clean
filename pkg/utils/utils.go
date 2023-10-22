package utils

import "os"

func GetEnvVar(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
