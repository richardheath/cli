package cli

import (
	"errors"
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

	commandArgs, flagArgs, err := splitCommandsAndFlags(args, flagPrefixes)
	if err != nil {
		return err
	}

	flagTypes := app.FlagTypes[:]
	command, err := getMatchingCommand(app.Commands, commandArgs, flagTypes)
	if err != nil {
		return err
	}

	knownFlags, unknownFlags, err := processFlags(flagArgs, flagTypes, flagPrefixes)
	err = command.Action(app, knownFlags, unknownFlags)
	if err != nil {
		return err
	}

	return nil
}

func splitCommandsAndFlags(args []string, flagPrefixes map[string]string) ([]string, []string, error) {
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

	return commands, flags, nil
}

func getMatchingCommand(commands []Command, commandArgs []string, types []FlagType) (Command, error) {
	for _, command := range commands {
		argPos := 0

		for _, path := range command.Path {
			isFlagBinder := strings.HasPrefix(path, "{") && strings.HasSuffix(path, "}")
			if commandArgs[argPos] == path || isFlagBinder {
				argPos++
			}
		}

		if argPos == len(commandArgs) {
			types = append(types, command.FlagTypes...)
			return command, nil
		}

		if argPos > 0 {
			types = append(types, command.FlagTypes...)
			return getMatchingCommand(command.Commands, commandArgs[argPos:], types)
		}
	}

	var noMatch Command
	return noMatch, errors.New("Command not found")
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
