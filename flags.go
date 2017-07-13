package cli

import "strings"
import "fmt"

// FlagPrefix Definition for flag groups. This can be used to specify different types of flag groups that the app supports.
type FlagPrefix struct {
	Key         string
	Shorthand   string
	Description string
}

// FlagType Flag type definition.
type FlagType struct {
	Key        string
	Shorthand  string
	Default    string
	Prefix     string
	Validators []FlagValidator
}

type ProcessedFlags struct {
	Known   map[string]string
	Unknown map[string]string
}

// FlagValidator Validator function that will be used to validate flag type value.
type FlagValidator func(value string) error

// FlagValidationFormatter Used to format message on validation error.
type FlagValidationFormatter func(flagType FlagType, value string, validationError error) string

// DefaultFlagValidationFormatter Default formatter.
func DefaultFlagValidationFormatter(flagType FlagType, value string, validationError error) string {
	return flagType.Key + ": " + value + "\n" + validationError.Error() + "\n"
}

// ProcessFlags Process arguments. Must call ProcessArgs first.
func (app *App) ProcessFlags(commandInfo CommandMatchInfo) (ProcessedFlags, error) {
	flagPrefixes := getFlagPrefixes(app)
	flags := ProcessedFlags{
		Known:   map[string]string{},
		Unknown: map[string]string{},
	}

	validatorErrors := ""
	flagArgs := append(app.FlagArgs, commandInfo.BindedFlags...)
	curFlagKey, curFlagPrefix, inFLagKey := "", "", false
	for _, arg := range flagArgs {
		prefix, hasPrefix := tryGetFlagPrefix(arg, flagPrefixes)

		if inFLagKey {
			value := arg
			flagType, flagTypeFound := getKeyFlagType(commandInfo.FlagTypes, flagPrefixes, curFlagKey, curFlagPrefix)

			// Set value to default when flag have no value.
			if hasPrefix {
				if flagTypeFound {
					value = flagType.Default
				} else {
					value = "true"
				}
			}

			if flagTypeFound {
				validationError := runFlagTypeValidators(flagType, value)
				if validationError != nil {
					validatorErrors += app.FlagValidationFormatter(flagType, value, validationError)
				}

				normalizedPrefix := flagPrefixes[curFlagPrefix]
				flags.Known[normalizedPrefix+flagType.Key] = value
			} else {
				flags.Unknown[curFlagPrefix+curFlagKey] = value
			}
		}

		if hasPrefix {
			curFlagKey = strings.Replace(arg, prefix, "", 1)
			curFlagPrefix = prefix
			inFLagKey = true
		} else {
			inFLagKey = false
		}
	}

	if validatorErrors != "" {
		return flags, fmt.Errorf(validatorErrors)
	}

	return flags, nil
}

func getKeyFlagType(flagTypes []FlagType, flagPrefixes map[string]string, flagKey string, flagPrefix string) (flagType FlagType, found bool) {
	var matchingType FlagType
	for _, flagType := range flagTypes {
		if (flagType.Key == flagKey || flagType.Shorthand == flagKey) && flagType.Prefix == flagPrefixes[flagPrefix] {
			return flagType, true
		}
	}

	return matchingType, false
}

func runFlagTypeValidators(flagType FlagType, flagValue string) error {
	errors := ""
	for _, validator := range flagType.Validators {
		err := validator(flagValue)
		if err != nil {
			errors += "  " + err.Error()
		}
	}

	if errors != "" {
		return fmt.Errorf(errors)
	}
	return nil
}

func tryGetFlagPrefix(value string, flagPrefixes map[string]string) (prefix string, hasPrefix bool) {
	for prefix := range flagPrefixes {
		if strings.HasPrefix(value, prefix) {
			return prefix, true
		}
	}
	return "", false
}

func insertFlagsDefaultValues(flags ProcessedFlags, commandInfo CommandMatchInfo) {
	for _, flag := range commandInfo.Command.FlagTypes {
		flagName := flag.Prefix + flag.Key

		if !isFlagKnown(flags, flagName) {
			flags.Known[flagName] = flag.Default
		}
	}
}

func isFlagKnown(flags ProcessedFlags, flagName string) bool {
	for key := range flags.Known {
		if key == flagName {
			return true
		}
	}

	return false
}
