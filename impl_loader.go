package confetti

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// implLoader implements ILoader.
type implLoader struct {
	opts *LoaderOptions
}

func (i *implLoader) Load(target interface{}) error {
	if !isStructPointer(target) {
		return errors.New("target must be a struct pointer")
	}

	// Fill out all options values.
	i.resolveOptions()

	// Getting the value out of the struct pointer to generate final JSON.
	structValue := reflect.ValueOf(target).Elem().Interface()
	result, err := i.generateJSON(structValue)
	if err != nil {
		return fmt.Errorf("failed to generate JSON: %w", err)
	}

	// Finally, unmarshalling into the provided target.
	if err := json.Unmarshal([]byte(result), target); err != nil {
		return fmt.Errorf("failed to unmarshal results into target: %w", err)
	}

	return nil
}

// generateJSON does TODO.
func (i *implLoader) generateJSON(structValue interface{}) (string, error) {
	var result string

	// Creating reflection types for looping over fields.
	reflectValue := reflect.ValueOf(structValue)
	reflectType := reflect.TypeOf(structValue)

	for ind := 0; ind < reflectValue.NumField(); ind++ {
		fieldReflectValue := reflectValue.Field(ind)
		fieldReflectType := reflectType.Field(ind)

		// If field is of pointer type, we extract the pointer's value.
		// This allows for cleaner code below.
		if fieldReflectValue.Kind() == reflect.Ptr {
			fieldReflectValue = fieldReflectValue.Elem()
		}

		// Resolving the value of the struct field.
		resolvedValue, err := i.resolveField(fieldReflectType)
		if err != nil {
			return "", fmt.Errorf("failed to resolve field: %w", err)
		}

		// If the resolved value is a string, it should be covered
		// in additional quotes to result in a valid JSON.
		switch asserted := resolvedValue.(type) {
		case string:
			resolvedValue = fmt.Sprintf(`"%s"`, asserted)
		}

		// If the field is of type struct, we will dive into it recursively.
		if fieldReflectValue.Kind() == reflect.Struct {
			// Loop over the struct fields again and merge the JSONs.
		}

		// Appending the results to the JSON string.
		result = fmt.Sprintf(`%s "%s":%s`, result, fieldReflectType.Name, resolvedValue)
	}

	// Checking for a trailing comma.
	if result[len(result)-1] == ',' {
		result = result[:len(result)-1]
	}

	return fmt.Sprintf("{%s}", result), nil
}

// resolveOptions does TODO.
func (i *implLoader) resolveOptions() {
	if i.opts == nil {
		i.opts = defaultLoaderOptions
		return
	}

	if i.opts.DefTagName == "" {
		i.opts.DefTagName = defaultLoaderOptions.DefTagName
	}
	if i.opts.EnvTagName == "" {
		i.opts.EnvTagName = defaultLoaderOptions.EnvTagName
	}
	if i.opts.ArgTagName == "" {
		i.opts.ArgTagName = defaultLoaderOptions.ArgTagName
	}
}

// resolveField does TODO.
func (i *implLoader) resolveField(field reflect.StructField) (interface{}, error) {
	return nil, nil
}
