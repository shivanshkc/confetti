package confetti

import (
	"flag"
	"reflect"
	"strings"
	"testing"
)

// TestImplFlagger_RegisterField tests if the RegisterField method correctly registers valid struct fields.
func TestImplFlagger_RegisterField(t *testing.T) {
	instance := &implFlagger{
		opts:    defaultLoaderOptions,
		flagSet: flag.NewFlagSet(defaultLoaderOptions.Title, flag.ContinueOnError),
		flags:   map[string]*customFlagHolder{},
	}

	dummyTarget := struct {
		// RegisterField should work with absent def tag.
		dummyField1 string `env:"DF1" arg:"df-1"`
		// RegisterField should work with absent env tag.
		dummyField2 string `def:"2" arg:"df-2,with doc"`
		// RegisterField should work with all tags present.
		dummyField3 string `def:"3" env:"DF3" arg:"df-3"`
	}{}

	var flags = []string{"df-1", "df-2", "df-3"}

	structValue := reflect.ValueOf(dummyTarget)
	structType := structValue.Type()

	for ind := 0; ind < structValue.NumField(); ind++ {
		fieldType := structType.Field(ind)
		// Registering the field.
		if err := instance.RegisterField(nil, &fieldType); err != nil {
			t.Errorf("Expected error: nil, got: %+v", err)
			return
		}
		// A registered flag must exist inside the flags map.
		if _, exists := instance.flags[flags[ind]]; !exists {
			t.Errorf("Expected flag %s to exist inside flags but it does not.", flags[ind])
			return
		}
	}
}

// TestImplFlagger_RegisterField_NoArgTag tests if the RegisterField method ignores fields with no arg tags.
func TestImplFlagger_RegisterField_NoArgTag(t *testing.T) {
	instance := &implFlagger{
		opts:    defaultLoaderOptions,
		flagSet: flag.NewFlagSet(defaultLoaderOptions.Title, flag.ContinueOnError),
		flags:   map[string]*customFlagHolder{},
	}

	dummyTarget := struct {
		// Arg tag completely absent.
		dummyField1 string `def:"1" env:"DF1"`
		// Arg tag empty.
		dummyField2 string `def:"2" env:"DF2" arg:""`
		// Arg tag non-empty but with only doc.
		dummyField3 string `def:"3" env:"DF3" arg:",only doc"`
	}{}

	structValue := reflect.ValueOf(dummyTarget)
	structType := structValue.Type()

	for ind := 0; ind < structValue.NumField(); ind++ {
		fieldType := structType.Field(ind)
		// Registering the field.
		if err := instance.RegisterField(nil, &fieldType); err != nil {
			t.Errorf("Expected error: nil, got: %+v", err)
			return
		}
	}
	// There should be no registrations.
	if len(instance.flags) != 0 {
		t.Errorf("Expected flags length to be 0, but got: %d", len(instance.flags))
		return
	}
}

// TestImplFlagger_RegisterField_Parse tests if the parse method works as expected with valid arguments.
func TestImplFlagger_RegisterField_Parse(t *testing.T) {
	instance := &implFlagger{
		opts:    defaultLoaderOptions,
		flagSet: flag.NewFlagSet(defaultLoaderOptions.Title, flag.ContinueOnError),
		flags:   map[string]*customFlagHolder{},
	}

	dummyTarget := struct {
		dummyField1 string `def:"1" env:"DF1" arg:"df-1,Dummy Field 1"`
		dummyField2 string `def:"2" env:"DF2" arg:"df-2,Dummy Field 2"`
	}{}

	structValue := reflect.ValueOf(dummyTarget)
	structType := structValue.Type()

	for ind := 0; ind < structValue.NumField(); ind++ {
		fieldType := structType.Field(ind)
		if err := instance.RegisterField(nil, &fieldType); err != nil {
			t.Errorf("Expected RegisterField error: nil, got: %+v", err)
			return
		}
	}

	// Ignoring unknown flags error.
	if err := instance.Parse(); err != nil && !strings.Contains(err.Error(), "flag provided but not defined") {
		t.Errorf("Expected Parse error: nil, got: %+v", err)
		return
	}
}

func TestImplFlagger_LookupFlag(t *testing.T) {
	instance := &implFlagger{
		opts:    defaultLoaderOptions,
		flagSet: flag.NewFlagSet(defaultLoaderOptions.Title, flag.ContinueOnError),
		flags:   map[string]*customFlagHolder{},
	}

	dummyTarget := struct {
		dummyField1 string `def:"1" env:"DF1" arg:"df-1"`
		dummyField2 string `def:"2" env:"DF2" arg:"df-2"`
		dummyField3 string `def:"3" env:"DF3" arg:"df-3"`
	}{}

	var flags = []string{"df-1", "df-2", "df-3"}

	structValue := reflect.ValueOf(dummyTarget)
	structType := structValue.Type()

	for ind := 0; ind < structValue.NumField(); ind++ {
		fieldType := structType.Field(ind)
		if err := instance.RegisterField(nil, &fieldType); err != nil {
			t.Errorf("Expected RegisterField error: nil, got: %+v", err)
			return
		}
	}

	// Ignoring unknown flags error.
	if err := instance.Parse(); err != nil && !strings.Contains(err.Error(), "flag provided but not defined") {
		t.Errorf("Expected Parse error: nil, got: %+v", err)
		return
	}

	// Testing if LookupFlag work as expected. Exists should be false because the flag is not being provided.
	for _, flg := range flags {
		if _, exists := instance.LookupFlag(flg); exists {
			t.Errorf("Expected flag: %s to not exist since it was not provided, but it exists.", flg)
			return
		}
	}
}
