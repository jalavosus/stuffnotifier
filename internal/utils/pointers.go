package utils

// ToPointer returns the pointer value
// of whatever is passed to it.
func ToPointer[T any](val T) *T {
	return &val
}

// FromPointer returns the dereferenced value
// of a pointer, along with a bool
// representing whether a nil value was passed,
// and thus, a zero value for the value's underlying
// type returned.
// If val is nil, the zero value for val's underlying
// type is returned, along with false.
func FromPointer[T any](val *T) (T, bool) {
	var fallback T

	if val != nil {
		return *val, true
	}

	return fallback, false
}
