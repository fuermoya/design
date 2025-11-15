package utils

func UPtr[T any](num T) *T {
	return &num
}
