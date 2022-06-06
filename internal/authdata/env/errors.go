package env

import (
	"fmt"
)

type KeyNotSetError struct {
	envKey string
}

func NewKeyNotSetError(envKey string) *KeyNotSetError {
	return &KeyNotSetError{envKey}
}

func (e *KeyNotSetError) Error() string {
	return fmt.Sprintf("key %[1]s not set in environment", e.envKey)
}
