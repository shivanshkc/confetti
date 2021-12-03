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
	split := strings.SplitN(tagValue, ",", 2)

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
		return nil, fmt.Errorf("failed to convert value: %w", err)
	}
	return converted, nil
}
