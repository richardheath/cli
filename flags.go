package cli

import "strings"

// FlagGroup Definition for flag groups. This can be used to specify different types of flag groups that the app supports.
type FlagGroup struct {
	Prefix          string
	ShorthandPrefix string
	Group           string
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

	flagKey := ""
	flagPrefix := ""
	for _, arg := range flagArgs {
		curPrefix := ""
		for prefix := range flagPrefixes {
			if strings.HasPrefix(arg, prefix) {
				curPrefix = prefix
				break
			}
		}

		if flagKey != "" {
			value := arg
			flagType, flagTypeFound := getKeyFlagType(flagTypes, flagPrefixes, flagKey, flagPrefix)

			// Set value to default when flag have no value.
			if curPrefix != "" {
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

				normalizedPrefix := flagPrefixes[flagPrefix]
				knownFlags[normalizedPrefix+flagType.Key] = value
			} else {
				unknownFlags[flagPrefix+flagKey] = value
			}
		}

		if curPrefix != "" {
			flagKey = strings.Replace(arg, curPrefix, "", 1)
			flagPrefix = curPrefix
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
