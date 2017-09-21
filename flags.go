package cli

// Flag Flag definition.
type Flag struct {
	Key         string
	Name        string
	Shorthand   string
	FlagType    string
	Description string
	Default     string
	Action      FlagAction
}

// FlagValues Storage for flag values.
type FlagValues map[string]string

// Flags Collection of flags.
type Flags []Flag

// FlagAction Method that will invoked when flag is found.
type FlagAction func(value string, ctx *Context) error

// GetValue Get flag value.
func (flags FlagValues) GetValue(flagType string, name string) (string, bool) {
	key := flagType + ":" + name
	for flagKey, flagValue := range flags {
		if flagKey == key {
			return flagValue, true
		}
	}

	return "", false
}

func (flags Flags) getFlag(key string) (Flag, bool) {
	for _, flag := range flags {
		if flag.Key == key {
			return flag, true
		}
	}

	return Flag{}, false
}
