package confetti

import "fmt"

var (
	// ErrNotStructPointer is returned when the provided target is not a struct pointer.
	ErrNotStructPointer = fmt.Errorf("target is not a struct pointer")
)
