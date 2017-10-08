package cli

import (
	"strings"
)

// Command Definition of an executable command.
type Command struct {
	Path     []string // Command path relative to its parent.
	Commands []*Command
	Flags    Flags
	Action   CommandAction
}

// CommandAction Method that will invoked when running a command.
type CommandAction func(ctx *Context) error

// Command Add sub command.
func (cmd *Command) Command(path string, action CommandAction) *Command {
	subCommand := Command{
		Path:     splitKeys(path),
		Action:   action,
		Flags:    Flags{},
		Commands: []*Command{},
	}

	cmd.Commands = append(cmd.Commands, &subCommand)
	return &subCommand
}

// Flag Add flag on current command.
func (cmd *Command) Flag(flagType string, name string, defaultValue string, action FlagAction) *Flag {
	keys := splitKeys(name)
	flag := Flag{
		Key:      flagType + ":" + keys[0],
		Name:     keys[0],
		FlagType: flagType,
		Default:  defaultValue,
		Action:   action,
	}

	if len(keys) > 1 {
		flag.Shorthand = keys[1]
	}

	cmd.Flags = append(cmd.Flags, flag)
	return &flag
}

func splitKeys(keys string) []string {
	return strings.Split(keys, " ")
}
