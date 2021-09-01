package confetti

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
)

const (
	tagArg     = "arg"
	tagEnv     = "env"
	tagDefault = "default"
)

// loaderImpl implements ILoader.
type loaderImpl struct {
	// cmdArgs contains all the passed command-line arguments. It is empty before the "Load" call.
	cmdArgs map[string]string
	// loaded contains the loaded configs. It is empty before the "Load" call.
	loaded map[string]interface{}
}

func (l *loaderImpl) Load(target interface{}) error {
	// The target should be a struct pointer.
	if !isStructPointer(target) {
		return ErrNotStructPointer
	}

	// Getting the value out of the struct pointer.
	structValue := reflect.ValueOf(target).Elem().Interface()

	// Loading command-line args.
	l.loadCmdArgs()

	l.loaded = map[string]interface{}{}
	// Looping over all struct fields and resolving the values according to tags.
	if err := forEachStructField(structValue, l.loadField, nil); err != nil {
		return err
	}

	// Marshalling the loaded values into JSON.
	loadedJSON, err := json.Marshal(l.loaded)
	if err != nil {
		return err
	}

	// Finally, unmarshalling the loaded values JSON into the target struct.
	return json.Unmarshal(loadedJSON, target)
}

// loadCmdArgs loops over and parses all the cmd args in os.Args[1:] according to `--arg-name=value` format.
func (l *loaderImpl) loadCmdArgs() {
	l.cmdArgs = map[string]string{}

	for _, arg := range os.Args {
		// Ignoring all args that do not begin with a double hyphen.
		if !strings.HasPrefix(arg, "--") {
			continue
		}

		// Removing the initial double hyphen.
		arg = arg[2:]
		if len(arg) == 0 {
			continue
		}

		// Splitting the arg by "=", the first element will be the name and the second will be the value.
		argNameAndValue := strings.SplitN(arg, "=", 2)
		if len(argNameAndValue) < 2 {
			continue
		}

		l.cmdArgs[argNameAndValue[0]] = argNameAndValue[1]
	}
}

// loadField loads the value of one struct field and adds its entry to the loaded map.
func (l *loaderImpl) loadField(parents []*reflect.StructField, field *reflect.StructField) error {
	resolvedValue, err := l.resolveFieldValue(field)
	if err != nil {
		return err
	}

	// Adding the nested entry inside the loaded map.
	finalValues := l.loaded
	for _, parent := range parents {
		if _, exists := finalValues[parent.Name]; !exists {
			finalValues[parent.Name] = map[string]interface{}{}
		}
		finalValues = finalValues[parent.Name].(map[string]interface{})
	}

	finalValues[field.Name] = resolvedValue
	return nil
}

// resolveFieldValue resolves the value of a struct field according to the specified struct tags.
func (l *loaderImpl) resolveFieldValue(field *reflect.StructField) (interface{}, error) {
	argName, present := field.Tag.Lookup(tagArg)
	if present {
		value, exists := l.cmdArgs[argName]
		if exists {
			return string2Interface(field.Type.Kind(), value)
		}
	}

	envName, present := field.Tag.Lookup(tagEnv)
	if present {
		value, exists := os.LookupEnv(envName)
		if exists {
			return string2Interface(field.Type.Kind(), value)
		}
	}

	defaultValue, exists := field.Tag.Lookup(tagDefault)
	if exists {
		return string2Interface(field.Type.Kind(), defaultValue)
	}

	return nil, nil
}

func (l *loaderImpl) init() {
	l.cmdArgs = map[string]string{}
	l.loaded = map[string]interface{}{}
}
