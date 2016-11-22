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

// FlagValidator Validator function that will be used to validate flag type value.
type FlagValidator func(value string) error

func processFlags(flagArgs []string, flagTypes []FlagType, flagPrefixes map[string]string) (knownFlags map[string]string, unknownFlags map[string]string, err error) {
	knownFlags = make(map[string]string)
	unknownFlags = make(map[string]string)
	validatorErrors := ""

	curFlagKey, curFlagPrefix, inFLagKey := "", "", false
	for _, arg := range flagArgs {
		prefix, hasPrefix := tryGetFlagPrefix(arg, flagPrefixes)

		if inFLagKey {
			value := arg
			flagType, flagTypeFound := getKeyFlagType(flagTypes, flagPrefixes, curFlagKey, curFlagPrefix)

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
					// TODO: Parameterize validation formatter.
					validatorErrors += flagType.Key + ": " + value + "\n" + validationError.Error() + "\n"
				}

				normalizedPrefix := flagPrefixes[curFlagPrefix]
				knownFlags[normalizedPrefix+flagType.Key] = value
			} else {
				unknownFlags[curFlagPrefix+curFlagKey] = value
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
		return knownFlags, unknownFlags, fmt.Errorf("Argument validation error:\n%s", validatorErrors)
	}

	return knownFlags, unknownFlags, nil
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
