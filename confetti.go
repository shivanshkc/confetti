package confetti

// ILoader represents a configuration loader.
type ILoader interface {
	// Load loads the configs into the provided struct.
	//
	// The target should be a struct pointer. The struct fields may have three tags:
	// "default", "env" and "arg".
	//
	// 1. The "default" tag will be accepted as the default value of the field.
	//
	// 2. The "env" tag should hold the name of environment variable bound to the field.
	//
	// 3. The "arg" tag should hold the name of the command-line flag bound to the field. Example:
	//    $ <script> --arg-name=value
	//
	// The order of priority is "arg" > "env" > "default".
	Load(target interface{}) error

	// init loads the implementation's dependencies. It panics on failure.
	init()
}

// GetLoader returns a new instance of the ILoader.
func GetLoader() ILoader {
	var loader ILoader = &loaderImpl{}
	loader.init()

	return loader
}
