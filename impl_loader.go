package confetti

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"

	"github.com/joho/godotenv"
)

// implLoader implements ILoader.
type implLoader struct {
	// opts keeps the LoaderOptions.
	opts *LoaderOptions
	// flagSet is responsible for parsing of flags.
	flagSet *flag.FlagSet
	// flags keeps track of all flag values.
	flags map[string]*customFlagHolder
}

func (i *implLoader) Load(target interface{}) error {
	// Validations.
	if !isStructPointer(target) {
		return errors.New("target must be a struct pointer")
	}

	// Fill out all options values.
	i.opts = mergeOptions(i.opts, defaultLoaderOptions)
	// Reading the .env file as per the option.
	if i.opts.UseDotEnv {
		_ = godotenv.Load()
	}

	// Getting the value out of the struct pointer to loop over its fields.
	structValue := reflect.ValueOf(target).Elem().Interface()

	// Parsing and persisting all required flags.
	i.setupFlags(structValue)

	// Generating the map for the struct.
	generatedMap, err := i.generateMap(structValue)
	if err != nil {
		return fmt.Errorf("failed to generate map: %w", err)
	}

	// Converting the generatedMap into JSON for unmarshalling.
	generatedJSON, err := json.Marshal(generatedMap)
	if err != nil {
		return fmt.Errorf("failed to marshal generated map: %+v", err)
	}

	// Finally, unmarshalling into the provided target.
	if err := json.Unmarshal(generatedJSON, target); err != nil {
		return fmt.Errorf("failed to unmarshal results into target: %w", err)
	}

	return nil
}

func (i *implLoader) setupFlags(structValue interface{}) {
	// Creating reflection types for looping over fields.
	reflectValue := reflect.ValueOf(structValue)
	reflectType := reflect.TypeOf(structValue)

	// If flagSet is nil, it means it is the first call.
	// Otherwise, it is a call through recursion.
	if i.flagSet == nil {
		// Initializing the flagSet. All further calls will be through recursion.
		i.flagSet = flag.NewFlagSet(reflectType.Name(), flag.ContinueOnError)
		i.flags = map[string]*customFlagHolder{}

		// The flags will be parsed when the function returns.
		defer func() { _ = i.flagSet.Parse(os.Args[1:]) }()
	}

	for ind := 0; ind < reflectValue.NumField(); ind++ {
		fieldReflectValue := reflectValue.Field(ind)
		fieldReflectType := reflectType.Field(ind)

		// Getting the value of the arg tag.
		argTagValue, present := fieldReflectType.Tag.Lookup(i.opts.ArgTagName)
		if !present || argTagValue == "" {
			continue
		}

		// The flagName is the name of the flag to be parsed.
		// The flagDoc is the usage info of the flag.
		flagName, flagDoc := getFlagNameAndDoc(argTagValue, ",")
		if flagName == "" {
			continue
		}
		if flagDoc != "" {
			// For a more understandable display.
			flagDoc += "\n"
		}

		// Using the def and env tag values to show even more info on "-h".
		defValue, present := fieldReflectType.Tag.Lookup(i.opts.DefTagName)
		if !present {
			defValue = "not provided"
		}

		envValue, present := fieldReflectType.Tag.Lookup(i.opts.EnvTagName)
		if !present || envValue == "" {
			envValue = "not provided"
		}

		// The usage instructions that will show up on "-h".
		usage := fmt.Sprintf("%sDefault: %s\nEnvironment: %s", flagDoc, defValue, envValue)

		// Binding the flag values to customFlagHolder.
		i.flags[flagName] = &customFlagHolder{}
		i.flagSet.Var(i.flags[flagName], flagName, usage)

		// If the field is of type struct, we will dive into it recursively.
		if fieldReflectValue.Kind() == reflect.Struct {
			i.setupFlags(fieldReflectValue.Interface())
		}
	}
}

// generateMap generates the map for the provided structValue according to the tags provided.
func (i *implLoader) generateMap(structValue interface{}) (interface{}, error) {
	generatedMap := map[string]interface{}{}

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

		// If the field is of type struct, we will dive into it recursively.
		if fieldReflectValue.Kind() == reflect.Struct {
			// This is the value of the struct as dictated by the tags on its fields.
			resolvedInsideValue, err := i.generateMap(fieldReflectValue.Interface())
			if err != nil {
				return "", err
			}

			// Populating inside value with any additional fields that outside value may have.
			mergeMaps(resolvedInsideValue.(map[string]interface{}), resolvedValue.(map[string]interface{}))
			// Ditching the outside value completely.
			resolvedValue = resolvedInsideValue
		}

		// Putting the resolved value in the map.
		generatedMap[fieldReflectType.Name] = resolvedValue
	}

	return generatedMap, nil
}

// resolveField uses the specified (or default) tag names to resolve value of a single struct field.
func (i *implLoader) resolveField(field reflect.StructField) (interface{}, error) {
	var tagValue string
	var present bool

	tagValue, present = field.Tag.Lookup(i.opts.ArgTagName)
	if present && tagValue != "" {
		// Getting only the flagName. We don't need flagDoc here.
		flagName, _ := getFlagNameAndDoc(tagValue, ",")

		holder, exists := i.flags[flagName]
		if exists && holder.setCalled {
			return string2Interface(field.Type.Kind(), holder.String())
		}
	}

	tagValue, present = field.Tag.Lookup(i.opts.EnvTagName)
	if present && tagValue != "" {
		value, exists := os.LookupEnv(tagValue)
		if exists {
			return string2Interface(field.Type.Kind(), value)
		}
	}

	tagValue, present = field.Tag.Lookup(i.opts.DefTagName)
	if present {
		return string2Interface(field.Type.Kind(), tagValue)
	}

	return nil, nil
}
