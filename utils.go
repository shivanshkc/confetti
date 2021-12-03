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

// mergeOptions looks for any absent fields in opts1 and fills them using opts2.
func mergeOptions(opts1 *LoaderOptions, opts2 *LoaderOptions) *LoaderOptions {
	if opts2 == nil {
		return opts1
	}
	if opts1 == nil {
		opts1 = opts2
		return opts1
	}

	if opts1.DefTagName == "" {
		opts1.DefTagName = opts2.DefTagName
	}
	if opts1.EnvTagName == "" {
		opts1.EnvTagName = opts2.EnvTagName
	}
	if opts1.ArgTagName == "" {
		opts1.ArgTagName = opts2.ArgTagName
	}

	return opts1
}

// mergeMaps populates map1 with any fields that are absent in it but present in map2.
// If both maps contain further nested maps, the merging happens recursively.
func mergeMaps(map1 map[string]interface{}, map2 map[string]interface{}) {
	for key, map2Value := range map2 {
		map1Value, exists := map1[key]
		if !exists {
			map1[key] = map2Value
			continue
		}

		subMap1, isSubMap1 := map1Value.(map[string]interface{})
		subMap2, isSubMap2 := map2Value.(map[string]interface{})
		if !isSubMap1 || !isSubMap2 {
			continue
		}

		mergeMaps(subMap1, subMap2)
	}
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
