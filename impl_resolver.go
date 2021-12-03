package confetti

import (
	"os"
)

// newResolver returns a new iResolver instance.
func newResolver(opts *LoaderOptions) iResolver {
	return &implResolver{opts: opts}
}

// implResolver implements iResolver.
type implResolver struct {
	// opts keeps the LoaderOptions.
	opts *LoaderOptions
}

func (i *implResolver) ResolveField(parents []rsf, field rsf, flagger iFlagger) (resolved interface{}, err error) {
	var tagValue string
	var present bool

	tagValue, present = field.Tag.Lookup(i.opts.ArgTagName)
	if present && tagValue != "" {
		// Getting only the flagName. We don't need flagDoc here.
		flagName, _ := getFlagNameAndDoc(tagValue, ",")
		flagValue, exists := flagger.LookupFlag(flagName)
		if exists {
			return string2Interface(field.Type.Kind(), flagValue)
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
