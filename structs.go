package confetti

// LoaderOptions can be used to customize the ILoader.
type LoaderOptions struct {
	// DefTagName can be used to alter the name of the def tag.
	DefTagName *string
	// EnvTagName can be used to alter the name of the env tag.
	EnvTagName *string
	// ArgTagName can be used to alter the name of the arg tag.
	ArgTagName *string
	// UseDotEnv controls whether to read data from the .env file.
	UseDotEnv *bool
}
