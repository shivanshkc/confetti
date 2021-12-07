package confetti

import (
	"errors"
	"reflect"
	"testing"
)

// implMockFlagger is a mock implementation of iFlagger.
type implMockFlagger struct {
	argMap      map[string]string
	registerErr error
}

func (i *implMockFlagger) RegisterField(_ []rsf, _ rsf) error { return i.registerErr }

func (i *implMockFlagger) Parse() error { return nil }

func (i *implMockFlagger) LookupFlag(flagName string) (flagValue string, exists bool) {
	value, exists := i.argMap[flagName]
	return value, exists
}

// implMockResolver is a mock implementation of iResolver.
type implMockResolver struct {
	errorMap map[string]error
	valueMap map[string]interface{}
}

func (i *implMockResolver) ResolveField(parents []rsf, field rsf, flagger iFlagger) (resolved interface{}, err error) {
	err, exists := i.errorMap[field.Name]
	if exists {
		return nil, err
	}
	value, exists := i.valueMap[field.Name]
	if exists {
		return value, nil
	}
	return nil, nil
}

// TestImplLoader_Load_NotStructPointer tests if the Load method gives
// an error when it is invoked without a struct pointer.
func TestImplLoader_Load_NotStructPointer(t *testing.T) {
	instance := &implLoader{}
	if err := instance.Load(2); err == nil {
		t.Errorf("Expected error from Load, but didn't get any.")
		return
	}
}

// TestImplLoader_Load tests if the Load method loads the correct values with valid input.
func TestImplLoader_Load(t *testing.T) {
	dummyTarget := struct {
		DummyField1 string
		DummyField2 struct {
			DummyField21 int
			DummyField22 struct {
				DummyField221 map[string]string
			}
		}
	}{}

	dummyField1Expected := "dummy1"
	dummyField21Expected := 10
	dummyField221Expected := map[string]string{"cool": "right"}

	instance := &implLoader{
		opts: defaultLoaderOptions,
		flagger: &implMockFlagger{
			argMap:      map[string]string{},
			registerErr: nil,
		},
		resolver: &implMockResolver{
			errorMap: map[string]error{},
			valueMap: map[string]interface{}{
				"DummyField1":   dummyField1Expected,
				"DummyField21":  dummyField21Expected,
				"DummyField221": dummyField221Expected,
			},
		},
	}

	if err := instance.Load(&dummyTarget); err != nil {
		t.Errorf("Expected error to be nil, but got: %+v", err)
		return
	}

	if dummyTarget.DummyField1 != dummyField1Expected {
		t.Errorf("Expected DummyField1 value to be: %+v, got: %+v", dummyField1Expected, dummyTarget.DummyField1)
	}
	if dummyTarget.DummyField2.DummyField21 != dummyField21Expected {
		t.Errorf("Expected DummyField31 value to be: %+v, got: %+v", dummyField21Expected, dummyTarget.DummyField2.DummyField21)
	}
	if !reflect.DeepEqual(dummyTarget.DummyField2.DummyField22.DummyField221, dummyField221Expected) {
		t.Errorf("Expected DummyField321 value to be: %+v, got: %+v", dummyField221Expected, dummyTarget.DummyField2.DummyField22.DummyField221)
	}
}

// TestImplLoader_Load_PointerError tests if the Load method gives an error on nested pointer usage.
func TestImplLoader_Load_PointerError(t *testing.T) {
	dummyTarget := struct {
		DummyField1 struct {
			DummyField11 *string
		}
	}{}

	instance := &implLoader{
		opts:     defaultLoaderOptions,
		flagger:  &implMockFlagger{},
		resolver: &implMockResolver{},
	}

	if err := instance.Load(&dummyTarget); err == nil {
		t.Errorf("Expected error from Load, but didn't get any.")
		return
	}
}

// TestImplLoader_Load_FlaggerError tests if the Load method gives an error if the Flagger fails.
func TestImplLoader_Load_FlaggerError(t *testing.T) {
	dummyTarget := struct {
		DummyField1 string
	}{}

	instance := &implLoader{
		opts:    defaultLoaderOptions,
		flagger: &implMockFlagger{registerErr: errors.New("failed to register")},
	}

	if err := instance.Load(&dummyTarget); err == nil {
		t.Errorf("Expected error from Load, but didn't get any.")
		return
	}
}

// TestImplLoader_Load_ResolverError tests if the Load method gives an error if the Resolver fails.
func TestImplLoader_Load_ResolverError(t *testing.T) {
	dummyTarget := struct {
		DummyField1 string
	}{}

	instance := &implLoader{
		opts:    defaultLoaderOptions,
		flagger: &implMockFlagger{},
		resolver: &implMockResolver{
			errorMap: map[string]error{"DummyField1": errors.New("failed to resolve field")},
		},
	}

	if err := instance.Load(&dummyTarget); err == nil {
		t.Errorf("Expected error from Load, but didn't get any.")
		return
	}
}
