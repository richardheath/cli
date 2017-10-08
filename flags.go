package cli

import (
	"strings"
)

// Flag Flag definition.
type Flag struct {
	Key       string
	Name      string
	Shorthand string
	FlagType  string
	Default   string
	Action    FlagAction
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

// SetValue Set flag value.
func (flags FlagValues) SetValue(flagType string, name string, value string) {
	key := flagType + ":" + name
	flags[key] = value
}

// GetFlagsByType Get all flags for given flag type.
func (flags FlagValues) GetFlagsByType(flagType string) []string {
	result := []string{}
	typePrefix := flagType + ":"
	for flagKey := range flags {
		if strings.HasPrefix(flagKey, typePrefix) {
			result = append(result, flagKey[len(typePrefix):])
		}
	}

	return result
}

func (flag Flag) nameKey() string {
	return flag.FlagType + ":" + flag.Name
}

func (flag Flag) shorthandKey() string {
	return flag.FlagType + ":" + flag.Shorthand
}

func (flags Flags) getFlag(flagType, flagName string) (Flag, bool) {
	for _, flag := range flags {
		flagMatches := flag.FlagType == flagType && (flag.Name == flagName || flag.Shorthand == flagName)
		if flagMatches {
			return flag, true
		}
	}

	return Flag{}, false
}

func splitFlagKey(key string) (flagType, flagName string) {
	split := strings.Split(key, ":")
	return split[0], split[1]
}
