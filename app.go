package cli

import (
	"fmt"
	"strings"
)

// App CLI Application
type App struct {
	Name         string
	Version      string
	FlagPrefixes []FlagPrefix
	FlagTypes    []FlagType
	Commands     []Command
}

// Run Execute command based on given arguments.
func (app App) Run(args []string) error {
	var err error
	flagPrefixes := getFlagPrefixes(app)
	commandArgs, flagArgs := splitCommandsAndFlags(args, flagPrefixes)

	command, flagTypes, bindedFlags, err := getMatchingCommand(app.Commands, commandArgs, app.FlagTypes, []string{})
	if err != nil {
		return err
	}

	flagArgs = append(flagArgs, bindedFlags...)
	knownFlags, unknownFlags, err := processFlags(flagArgs, flagTypes, flagPrefixes)
	err = command.Action(app, knownFlags, unknownFlags)
	if err != nil {
		return err
	}

	return nil
}

func splitCommandsAndFlags(args []string, flagPrefixes map[string]string) ([]string, []string) {
	commands := make([]string, 0, 2)
	flags := make([]string, 0, 2)
	var currentFlag string

	for _, arg := range args {
		for prefix := range flagPrefixes {
			if strings.HasPrefix(arg, prefix) {
				currentFlag = arg
				break
			}
		}

		if len(currentFlag) == 0 {
			commands = append(commands, arg)
		} else {
			flags = append(flags, arg)

			if currentFlag != arg {
				currentFlag = ""
			}
		}
	}

	return commands, flags
}

func getMatchingCommand(commands []Command, commandArgs []string, types []FlagType, bindedFlags []string) (Command, []FlagType, []string, error) {
	for _, command := range commands {
		argPos := 0

		for _, path := range command.Path {
			isFlagBinder := strings.HasPrefix(path, "{{") && strings.HasSuffix(path, "}}")
			if isFlagBinder {
				flagKey := path[2 : len(path)-2]

				bindedFlags = append(bindedFlags, flagKey, commandArgs[argPos])
			} else {
			}

			if commandArgs[argPos] == path || isFlagBinder {
				argPos++
			}
		}

		if argPos == len(commandArgs) {
			types = append(types, command.FlagTypes...)
			return command, types, bindedFlags, nil
		}

		if argPos > 0 {
			types = append(types, command.FlagTypes...)
			return getMatchingCommand(command.Commands, commandArgs[argPos:], types, bindedFlags)
		}
	}

	var noMatch Command
	return noMatch, types, bindedFlags, fmt.Errorf("Command not found: %v", commandArgs)
}

func getFlagPrefixes(app App) map[string]string {
	prefixes := make(map[string]string)
	for _, prefix := range app.FlagPrefixes {
		if prefix.Key != "" {
			prefixes[prefix.Key] = prefix.Key
		}

		if prefix.Shorthand != "" {
			prefixes[prefix.Shorthand] = prefix.Key
		}
	}

	return prefixes
}
