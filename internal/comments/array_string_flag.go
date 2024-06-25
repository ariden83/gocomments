package comments

import (
	"fmt"
)

// ArrayStringFlag are defined for string flags that may have multiple values.
type ArrayStringFlag []string

// Returns the concatenated string representation of the array of flags.
func (f *ArrayStringFlag) String() string {
	return fmt.Sprintf("%v", *f)
}

// Get returns an empty interface that may be type-asserted to the underlying
// value of type bool, string, etc.
func (f *ArrayStringFlag) Get() interface{} {
	return ""
}

// Set appends value the array of flags.
func (f *ArrayStringFlag) Set(value string) error {
	*f = append(*f, value)
	return nil
}
