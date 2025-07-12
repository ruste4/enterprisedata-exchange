package lib

import "github.com/google/uuid"

func GenerateUUID() string {
	return uuid.New().String()
}

func GenerateUUIDWithPrefix(prefix string) string {
	return prefix + "-" + GenerateUUID()
}
