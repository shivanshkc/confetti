package confetti

import (
	"fmt"
	"os"
)

// implResolver implements iResolver.
type implResolver struct {
	// opts keeps the LoaderOptions.
	opts *LoaderOptions
}

func (i *implResolver) ResolveField(parents []rsf, field rsf, flagger iFlagger) (interface{}, error) {
	resolveErr := fmt.Errorf(`failed to resolve field: "%s"`, formatNestedFieldName(parents, field))

	stringValue, present := i.resolveArg(field, flagger)
	if present {
		value, err := string2Interface(field.Type.Kind(), stringValue)
		return value, checkAndWrapErr(err, resolveErr)
	}

	stringValue, present = i.resolveEnv(field)
	if present {
		value, err := string2Interface(field.Type.Kind(), stringValue)
		return value, checkAndWrapErr(err, resolveErr)
	}

	stringValue, present = i.resolveDef(field)
	if present {
		value, err := string2Interface(field.Type.Kind(), stringValue)
		return value, checkAndWrapErr(err, resolveErr)
	}

	return nil, nil
}

func (i *implResolver) resolveArg(field rsf, flagger iFlagger) (string, bool) {
	tagValue, present := field.Tag.Lookup(i.opts.ArgTagName)
	if !present || tagValue == "" {
		return "", false
	}

	// Getting only the flagName. We don't need flagDoc here.
	flagName, _ := getFlagNameAndDoc(tagValue, ",")
	return flagger.LookupFlag(flagName)
}

func (i *implResolver) resolveEnv(field rsf) (string, bool) {
	tagValue, present := field.Tag.Lookup(i.opts.EnvTagName)
	if !present || tagValue == "" {
		return "", false
	}

	return os.LookupEnv(tagValue)
}

func (i *implResolver) resolveDef(field rsf) (string, bool) {
	return field.Tag.Lookup(i.opts.DefTagName)
}
