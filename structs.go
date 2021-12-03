package confetti

// defaultLoaderOptions are used when the user does not provide any.
var defaultLoaderOptions = &LoaderOptions{
	DefTagName: "def",
	EnvTagName: "env",
	ArgTagName: "arg",
	UseDotEnv:  false,
}

// LoaderOptions can be used to customize the ILoader.
type LoaderOptions struct {
	// DefTagName can be used to alter the name of the def tag.
	DefTagName string
	// EnvTagName can be used to alter the name of the env tag.
	EnvTagName string
	// ArgTagName can be used to alter the name of the arg tag.
	ArgTagName string
	// UseDotEnv controls whether to read data from the .env file.
	UseDotEnv bool
}

// complete checks all fields in the struct and fills in any absent ones using the default options.
func (l *LoaderOptions) complete() {
	if l.DefTagName == "" {
		l.DefTagName = defaultLoaderOptions.DefTagName
	}
	if l.EnvTagName == "" {
		l.EnvTagName = defaultLoaderOptions.EnvTagName
	}
	if l.ArgTagName == "" {
		l.ArgTagName = defaultLoaderOptions.ArgTagName
	}
}

// customFlagHolder keeps track of the flagValue, and whether it was ever set or not.
type customFlagHolder struct {
	// flagValue is the value of the flag.
	flagValue string
	// exists is true only if the flagValue has been set at least once.
	exists bool
}

func (c *customFlagHolder) String() string {
	return c.flagValue
}

func (c *customFlagHolder) Set(s string) error {
	c.exists = true
	c.flagValue = s
	return nil
}
