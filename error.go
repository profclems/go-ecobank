package ecobank

import "strings"

// ResponseError represents a collection of error messages.
type ResponseError []string

// Add adds a new error message to the collection.
func (e *ResponseError) Add(err string) {
	*e = append(*e, err)
}

// Error returns a formatted string of all error messages. Implementing the error interface.
func (e *ResponseError) Error() string {
	if e.Len() == 0 {
		return ""
	}
	return strings.Join(*e, "\n")
}

// All returns all error messages.
func (e *ResponseError) All() []string {
	if e.Len() == 0 {
		return nil
	}

	return *e
}

// Len returns the number of error messages.
func (e *ResponseError) Len() int {
	return len(*e)
}

// String returns a formatted string of all error messages.
func (e *ResponseError) String() string {
	return e.Error()
}
