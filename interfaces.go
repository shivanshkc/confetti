package confetti

// ILoader represents a configuration loader.
type ILoader interface {
	// Load loads the configs into the provided target.
	Load(target interface{}) error
}

// NewLoader provides a new ILoader instance.
func NewLoader(opts *LoaderOptions) ILoader {
	return &implLoader{opts: opts}
}
