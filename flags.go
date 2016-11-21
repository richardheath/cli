package cli

import "strings"

// FlagPrefix Definition for flag groups. This can be used to specify different types of flag groups that the app supports.
type FlagPrefix struct {
	Key         string
	Shorthand   string
	Description string
}

type FlagGroupValues struct {
	Known   map[string]string
	Unknown map[string]string
}

type FlagType struct {
	Key        string
	Shorthand  string
	Default    string
	Group      string
	Validators []FlagValidator
}

type FlagValidator func(key string, value string) error

func processFlags(flagArgs []string, flagTypes []FlagType, flagPrefixes map[string]string) (knownFlags map[string]string, unknownFlags map[string]string, err error) {
	knownFlags = make(map[string]string)
	unknownFlags = make(map[string]string)
	validatorErrors := ""

	curFlagKey := ""
	curFlagPrefix := ""
	for _, arg := range flagArgs {
		prefix, hasPrefix := tryGetFlagPrefix(arg, flagPrefixes)

		if curFlagKey != "" {
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
					validatorErrors += flagType.Key + " " + flagType.Group + " validation error\n" + validationError.Error()
				}

				normalizedPrefix := flagPrefixes[curFlagPrefix]
				knownFlags[normalizedPrefix+flagType.Key] = value
			} else {
				unknownFlags[curFlagPrefix+curFlagKey] = value
			}
		}

		if prefix != "" {
			curFlagKey = strings.Replace(arg, prefix, "", 1)
			curFlagPrefix = prefix
		}
	}

	return knownFlags, unknownFlags, nil
}

func getKeyFlagType(flagTypes []FlagType, flagPrefixes map[string]string, flagKey string, flagPrefix string) (flagType FlagType, found bool) {
	var matchingType FlagType
	for _, flagType := range flagTypes {
		if (flagType.Key == flagKey || flagType.Shorthand == flagKey) && flagType.Group == flagPrefixes[flagPrefix] {
			return flagType, true
		}
	}

	return matchingType, false
}

func runFlagTypeValidators(flagType FlagType, flagValue string) error {
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
