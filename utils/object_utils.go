package utils

// Stringable represents a type that can be converted into a string
type Stringable interface {
	// ToString returns a string representation of the type
	ToString() string
}
