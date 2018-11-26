package utils

type Comparable interface {
	CompareTo(anotherObj Comparable) bool
}

type Stringable interface {
	ToString() string
}