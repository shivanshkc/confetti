package confetti

import (
	"flag"
)

// ILoader represents a configuration loader.
type ILoader interface {
	// Load loads the configs into the provided target.
	Load(target interface{}) error
}

// iFlagger manages the flag parsing and persistence.
type iFlagger interface {
	// RegisterField registers a struct field into the underlying flagSet using the various struct tags.
	RegisterField(parents []rsf, field rsf) error
	// Parse parses the flags. It should be called after all RegisterField calls and before all LookupFlag calls.
	Parse() error
	// LookupFlag provides the value of the specified flag. The second return param tells if the value exists.
	LookupFlag(flagName string) (flagValue string, exists bool)
}

// iResolver manages the resolution of values.
type iResolver interface {
	// ResolveField resolves the value of a struct field using the various struct tags.
	// It requires an iFlagger to get the flag values.
	ResolveField(parents []rsf, field rsf, flagger iFlagger) (resolved interface{}, err error)
}

// NewDefLoader provides a new ILoader instance with default settings.
func NewDefLoader() ILoader {
	return NewLoader(*defaultLoaderOptions)
}

// NewLoader provides a new ILoader instance.
func NewLoader(opts LoaderOptions) ILoader {
	// Filling out missing option values.
	opts.complete()
	return &implLoader{
		opts:     &opts,
		flagger:  newFlagger(&opts),
		resolver: newResolver(&opts),
	}
}

// newFlagger returns a new iFlagger instance.
func newFlagger(opts *LoaderOptions) iFlagger {
	return &implFlagger{
		opts:    opts,
		flagSet: flag.NewFlagSet(opts.Title, flag.ContinueOnError),
		flags:   map[string]*customFlagHolder{},
	}
}

// newResolver returns a new iResolver instance.
func newResolver(opts *LoaderOptions) iResolver {
	return &implResolver{opts: opts}
}
