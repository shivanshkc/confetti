package confetti

import (
	"flag"
	"fmt"
	"os"
)

// implFlagger implements iFlagger.
type implFlagger struct {
	// opts keeps the LoaderOptions.
	opts *LoaderOptions
	// flagSet is responsible for parsing of flags.
	flagSet *flag.FlagSet
	// flags keeps track of all flag values.
	flags map[string]*customFlagHolder
}

func (i *implFlagger) Parse() error {
	if err := i.flagSet.Parse(os.Args[1:]); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}
	return nil
}

func (i *implFlagger) RegisterField(parents []rsf, field rsf) error {
	// Getting the value of the arg tag.
	argTagValue, present := field.Tag.Lookup(i.opts.ArgTagName)
	if !present || argTagValue == "" {
		return nil
	}

	// The flagName is the name of the flag to be parsed.
	// The flagDoc is the usage info of the flag.
	flagName, flagDoc := getFlagNameAndDoc(argTagValue, ",")
	if flagName == "" {
		return nil
	}
	if flagDoc != "" {
		// For a more understandable display.
		flagDoc += "\n"
	}

	// Using the def and env tag values to show even more info on "-h".
	defValue, present := field.Tag.Lookup(i.opts.DefTagName)
	if !present {
		defValue = "not provided"
	}

	envValue, present := field.Tag.Lookup(i.opts.EnvTagName)
	if !present || envValue == "" {
		envValue = "not provided"
	}

	// The usage instructions that will show up on "-h".
	usage := fmt.Sprintf("%sDefault: %s\nEnvironment: %s", flagDoc, defValue, envValue)

	// Binding the flag values to customFlagHolder.
	i.flags[flagName] = &customFlagHolder{}
	i.flagSet.Var(i.flags[flagName], flagName, usage)

	return nil
}

func (i *implFlagger) LookupFlag(flagName string) (flagValue string, exists bool) {
	holder, exists := i.flags[flagName]
	return holder.String(), exists && holder.exists
}
