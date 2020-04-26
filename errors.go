package libuconf

import (
	"errors"
	"fmt"
)

// Base errors that the package might return.
var (
	ErrNoVal = errors.New("missing value")
	ErrSet   = errors.New("failed to set")
)

func errNoVal(name string) error {
	return fmt.Errorf("%w: %s", ErrNoVal, name)
}
func errSet(name, val string) error {
	return fmt.Errorf("%w: %s to %s", ErrSet, name, val)
}
