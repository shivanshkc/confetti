package confetti

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/joho/godotenv"
)

// implLoader implements ILoader.
type implLoader struct {
	// opts keeps the LoaderOptions.
	opts *LoaderOptions

	flagger  iFlagger
	resolver iResolver
}

func (i *implLoader) Load(target interface{}) error {
	// Validations.
	if !isStructPointer(target) {
		return errors.New("target must be a struct pointer")
	}

	// Reading the .env file as per the option.
	if i.opts.UseDotEnv {
		_ = godotenv.Load()
	}

	// Getting the value out of the struct pointer to loop over its fields.
	structValue := reflect.ValueOf(target).Elem().Interface()

	// Creating the flagSet.
	if err := i.forEachStructField(structValue, i.flagger.RegisterField, nil); err != nil {
		return fmt.Errorf("failed to create flagSet: %w", err)
	}

	// Parsing all the flags.
	if err := i.flagger.Parse(); err != nil {
		return err
	}

	targetMap := msi{}
	// Loading all values inside the targetMap.
	if err := i.forEachStructField(structValue, i.resolveFieldWrapper(targetMap), nil); err != nil {
		return fmt.Errorf("failed to resolve values: %w", err)
	}

	// Marshalling the targetMap values into JSON.
	targetJSON, err := json.Marshal(targetMap)
	if err != nil {
		return fmt.Errorf("failed to marshal the targetMap: %w", err)
	}

	// Finally, unmarshalling the JSON into the target struct.
	if err := json.Unmarshal(targetJSON, target); err != nil {
		return fmt.Errorf("failed to unmarshal into target: %w", err)
	}

	return nil
}

// forEachStructField loops over all the fields of the provided input (struct)
// and calls action for each of those fields.
//
// The "value" argument should be a struct (not even a struct pointer).
//
// The "action" argument is the function to be called for each field.
//
// The "parents" argument is received by recursive calls made internally.
// External calls should provide this value as nil.
func (i *implLoader) forEachStructField(value interface{}, action structFieldAction, parents []rsf) error {
	// Creating reflection types for looping over fields.
	reflectValue := reflect.ValueOf(value)
	reflectType := reflect.TypeOf(value)

	// Looping over all struct fields.
	for ind := 0; ind < reflectValue.NumField(); ind++ {
		fieldReflectValue := reflectValue.Field(ind)
		fieldReflectType := reflectType.Field(ind)

		// If field is of pointer type, we extract the pointer's value. This allows for cleaner code below.
		if fieldReflectValue.Kind() == reflect.Ptr {
			fieldReflectValue = fieldReflectValue.Elem()
		}

		// Calling action of the struct field.
		if err := action(parents, &fieldReflectType); err != nil {
			return err
		}

		// If the field is not a struct, there's nothing else to be done.
		if fieldReflectValue.Kind() != reflect.Struct {
			continue
		}

		// Appending the parent before doing the recursive call.
		newParents := append(parents, &fieldReflectType)
		// Looping over the nested struct fields.
		if err := i.forEachStructField(fieldReflectValue.Interface(), action, newParents); err != nil {
			return err
		}
	}

	return nil
}

// resolveFieldWrapper is a wrapper around the iResolver.ResolveField method to
// make it a valid structFieldAction while also putting the targetMap in the scope.
func (i *implLoader) resolveFieldWrapper(targetMap msi) structFieldAction {
	return func(parents []rsf, field rsf) error {
		// Getting the resolved value.
		resolved, err := i.resolver.ResolveField(parents, field, i.flagger)
		if err != nil {
			return err
		}

		// Creating a separate map.
		// So, we won't lose the reference to the original targetMap.
		nestedMap := targetMap
		// Creating a nested entry inside the targetMap.
		for _, parent := range parents {
			if _, exists := nestedMap[parent.Name]; !exists {
				nestedMap[parent.Name] = msi{}
			}
			nestedMap = nestedMap[parent.Name].(msi)
		}

		// If the field type is struct, then the resolved value should be an empty map instead of nil.
		if resolved == nil && field.Type.Kind() == reflect.Struct {
			resolved = msi{}
		}
		nestedMap[field.Name] = resolved
		return nil
	}
}
