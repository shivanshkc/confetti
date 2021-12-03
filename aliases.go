package confetti

import (
	"reflect"
)

// rsf is a type alias for *reflect.StructField
type rsf = *reflect.StructField

// msi is a type alias for map[string]interface{}
type msi = map[string]interface{}

// structFieldAction represents an action that can be performed over a nested struct field.
// It handles nested fields well because it receives all the parents of the field as well.
type structFieldAction func(parents []rsf, field rsf) error
