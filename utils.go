package confetti

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// isStructPointer returns true if the input is a struct pointer, otherwise false.
func isStructPointer(input interface{}) bool {
	value := reflect.ValueOf(input)
	return value.Kind() == reflect.Ptr && value.Elem().Kind() == reflect.Struct
}

// getFlagNameAndDoc accepts the tag value of the "arg" tag and returns the flagName and flagDoc.
// The function assumes that the name and doc are separated by the provided separator.
func getFlagNameAndDoc(tagValue string, sep string) (flagName string, flagDoc string) {
	split := strings.SplitN(tagValue, sep, 2)

	switch len(split) {
	case 0:
		return "", ""
	case 1:
		return split[0], ""
	case 2:
		return split[0], split[1]
	default:
		return "", ""
	}
}

// string2Interface converts string values to JSON.
//
// Note that int, float, booleans etc. are also valid JSON.
func string2Interface(kind reflect.Kind, value string) (interface{}, error) {
	if kind == reflect.String {
		return value, nil
	}

	var converted interface{}
	if err := json.Unmarshal([]byte(value), &converted); err != nil {
		return nil, err
	}
	return converted, nil
}

// formatNestedFieldName accepts a field and its parents to create a formatted name string.
// Example: Parent1.Parent2.MyField
func formatNestedFieldName(parents []rsf, field rsf) string {
	var formatted string
	for _, parent := range parents {
		formatted += fmt.Sprintf("%s.", parent.Name)
	}
	formatted += field.Name
	return formatted
}

// checkAndWrapErr returns nil if 'err' is nil.
// If 'err' is not nil, it wraps the 'err' in 'wrappingErr' and returns it.
func checkAndWrapErr(err error, wrappingErr error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s %w", wrappingErr.Error(), err)
}
