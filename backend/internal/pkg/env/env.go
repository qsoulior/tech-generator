package env

import (
	"os"
	"strconv"
)

func GetString(key string, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return value
}

func GetInteger(key string, fallback int) int {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsedValue
}

func GetFloat(key string, fallback float64) float64 {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	parsedValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fallback
	}
	return parsedValue
}
