package confetti

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// structFieldAction represents an action that can be performed over a nested struct field.
// It handles nested fields well because it receives all the parents of the field as well.
type structFieldAction func(parents []*reflect.StructField, field *reflect.StructField) error

// isStructPointer returns true if the input is a struct pointer, otherwise false.
func isStructPointer(input interface{}) bool {
	value := reflect.ValueOf(input)
	if value.Kind() != reflect.Ptr {
		return false
	}

	return value.Elem().Kind() == reflect.Struct
}

// forEachStructField loops over all the fields of the provided input (struct)
// and calls action for each of those fields.
//
// The "input" argument should be a struct (not even a struct pointer).
//
// The "action" argument is the function to be called for each field.
//
// The "parents" argument is received by recursive calls made internally.
// In most cases, user should provide this value as nil.
func forEachStructField(input interface{}, action structFieldAction, parents []*reflect.StructField) error {
	inputValue := reflect.ValueOf(input)
	inputType := reflect.TypeOf(input)

	// This loop runs N times, where N is the number of fields in the input.
	for i := 0; i < inputValue.NumField(); i++ {
		fieldValue := inputValue.Field(i)
		fieldType := inputType.Field(i)

		// If value is a pointer, then we extract the pointer's value. This allows for cleaner code below.
		if fieldValue.Kind() == reflect.Ptr {
			fieldValue = fieldValue.Elem()
		}

		switch fieldValue.Kind() {
		// If the field is of type struct, we recursively loop over its fields.
		case reflect.Struct:
			newParents := append(parents, &fieldType)
			if err := forEachStructField(fieldValue.Interface(), action, newParents); err != nil {
				return err
			}
		// If the field is other than a struct, we execute the action for it.
		default:
			if err := action(parents, &fieldType); err != nil {
				return err
			}
		}
	}

	return nil
}

// string2Interface converts string values to JSON.
//
// Note that ints, floats, booleans etc are also valid JSON.
func string2Interface(kind reflect.Kind, value string) (interface{}, error) {
	if kind == reflect.String {
		return value, nil
	}

	var converted interface{}
	if err := json.Unmarshal([]byte(value), &converted); err != nil {
		return nil, fmt.Errorf("confetti: failed to convert value: %w", err)
	}

	return converted, nil
}
