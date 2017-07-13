package cli

import (
	"fmt"
	"strings"
)

// Command Definition of an executable command.
// Commands are matched based on their path.
type Command struct {
	Path      []string // Command path relative to its parent.
	Commands  []Command
	FlagTypes []FlagType
	Usage     string
	Action    CommandAction
}

// CommandMatchInfo Metadata for matching command. This includes
type CommandMatchInfo struct {
	unprocessedArgs []string
	BindedFlags     []string
	FlagTypes       []FlagType
	Command         Command
}

// CommandAction Method that will invoked when running a command.
type CommandAction func(flags ProcessedFlags) error

func getMatchingCommand(commands []Command, match *CommandMatchInfo) error {
	for _, command := range commands {
		argPos := 0

		for _, path := range command.Path {
			isFlagBinder := strings.HasPrefix(path, "{{") && strings.HasSuffix(path, "}}")
			if isFlagBinder {
				flagKey := path[2 : len(path)-2]
				match.BindedFlags = append(match.BindedFlags, flagKey, match.unprocessedArgs[argPos])
			}

			if match.unprocessedArgs[argPos] == path || isFlagBinder {
				argPos++
			}
		}

		if argPos == len(match.unprocessedArgs) {
			match.FlagTypes = append(match.FlagTypes, command.FlagTypes...)
			match.Command = command
			return nil
		}

		if argPos > 0 {
			match.FlagTypes = append(match.FlagTypes, command.FlagTypes...)
			match.unprocessedArgs = match.unprocessedArgs[argPos:]
			return getMatchingCommand(command.Commands, match)
		}
	}

	return fmt.Errorf("Command not found")
}
