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

	// Getting the value out of the struct pointer.
	structValue := reflect.ValueOf(target).Elem().Interface()

	// Creating reflection types for looping over fields.
	reflectValue := reflect.ValueOf(structValue)
	reflectType := reflect.TypeOf(structValue)

	// The resolved field values will be stored in this map.
	results := map[string]interface{}{}
	for i := 0; i < reflectValue.NumField(); i++ {
		fmt.Println("Field:", reflectType.Field(i))
	}

	// Marshalling the results map into JSON for later unmarshalling.
	resultJSON, err := json.Marshal(results)
	if err != nil {
		return fmt.Errorf("failed to marshal results: %w", err)
	}

	// Finally, unmarshalling into the provided target.
	if err := json.Unmarshal(resultJSON, target); err != nil {
		return fmt.Errorf("failed to unmarshal results into target: %w", err)
	}

	return nil
}
