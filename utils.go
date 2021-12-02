package confetti

import (
	"reflect"
)

// isStructPointer returns true if the input is a struct pointer, otherwise false.
func isStructPointer(input interface{}) bool {
	value := reflect.ValueOf(input)
	return value.Kind() == reflect.Ptr && value.Elem().Kind() == reflect.Struct
}
