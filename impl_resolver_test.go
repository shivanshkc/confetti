package confetti

import (
	"os"
	"reflect"
	"testing"
)

// implMockFlagger is a mock implementation of iFlagger.
type implMockFlagger struct {
	argMap map[string]string
}

func (i *implMockFlagger) RegisterField(_ []rsf, _ rsf) error { return nil }

func (i *implMockFlagger) Parse() error { return nil }

func (i *implMockFlagger) LookupFlag(flagName string) (flagValue string, exists bool) {
	value, exists := i.argMap[flagName]
	return value, exists
}

// TestImplResolver_ResolveField tests if the ResolveField method works as expected with correct inputs.
func TestImplResolver_ResolveField(t *testing.T) {
	instance := &implResolver{opts: defaultLoaderOptions}
	flagger := &implMockFlagger{argMap: map[string]string{}}

	// The expected resolved value of each field is governed by the mockers slice below.
	dummyTarget := struct {
		dummyField1 string `def:"1" env:"" arg:""`
		dummyField2 string `def:"2" env:"DF2" arg:"df-2"`
		dummyField3 string `def:"3" env:"DF3" arg:"df-3"`
		dummyField4 string `env:"DF4" arg:"df-4"`
	}{}

	mockers := []func() interface{}{
		func() interface{} {
			// Def tag will take effect.
			return "1"
		},
		func() interface{} {
			// Env tag will take effect.
			value := "env-value"
			_ = os.Setenv("DF2", value)
			return value
		},
		func() interface{} {
			// Arg tag will take effect.
			value := "arg-value"
			flagger.argMap["df-3"] = value
			return value
		},
		func() interface{} {
			// No applicable value.
			return nil
		},
	}

	structValue := reflect.ValueOf(dummyTarget)
	structType := structValue.Type()

	for ind := 0; ind < structValue.NumField(); ind++ {
		fieldType := structType.Field(ind)

		// Mocking the field value.
		expected := mockers[ind]()

		// Resolving the field value.
		resolved, err := instance.ResolveField(nil, &fieldType, flagger)
		if err != nil {
			t.Errorf("Expecting no error in ResolveField, but got: %+v", err)
			return
		}

		// Verifying if ResolveField got the correct value.
		if expected != resolved {
			t.Errorf("expected resolved value: %+v, but got: %+v", expected, resolved)
			return
		}
	}
}

// TestImplResolver_ResolveField_BadType tests if ResolveField returns an error upon bad value types.
func TestImplResolver_ResolveField_BadType(t *testing.T) {
	instance := &implResolver{opts: defaultLoaderOptions}
	flagger := &implMockFlagger{argMap: map[string]string{}}

	// The expected resolved value of each field is governed by the mockers slice below.
	// But note that the values are of bad data types, so they will not be resolved.
	dummyTarget := struct {
		dummyField1 int     `def:"{{" env:"" arg:""`
		dummyField2 float32 `def:"32.1" env:"DF2" arg:"df-2"`
		dummyField3 bool    `def:"true" env:"DF3" arg:"df-3"`
	}{}

	mockers := []func() interface{}{
		func() interface{} {
			// Def tag will take effect.
			return "1"
		},
		func() interface{} {
			// Env tag will take effect.
			value := "{{"
			_ = os.Setenv("DF2", value)
			return value
		},
		func() interface{} {
			// Arg tag will take effect.
			value := "{{"
			flagger.argMap["df-3"] = value
			return value
		},
	}

	structValue := reflect.ValueOf(dummyTarget)
	structType := structValue.Type()

	for ind := 0; ind < structValue.NumField(); ind++ {
		fieldType := structType.Field(ind)

		// Mocking the field value.
		mockers[ind]()

		// Resolving the field value.
		resolved, err := instance.ResolveField(nil, &fieldType, flagger)
		if err == nil {
			t.Errorf("expected err to occur but got resolved value: %+v", resolved)
			return
		}
	}
}
