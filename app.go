package cli

import (
	"fmt"
	"strings"
)

// App CLI Application
type App struct {
	Name         string
	Version      string
	baseCommand  *Command
	flagPrefixes map[string]string
}

// Context Context of CLI execution.
type Context struct {
	App   *App
	Flags FlagValues
}

// NewApp Initialize App struct.
func NewApp(name string, version string) App {
	return App{
		Name:    name,
		Version: version,
		baseCommand: &Command{
			Path:     []string{"base"},
			Commands: []*Command{},
		},
		flagPrefixes: map[string]string{},
	}
}

// Run Execute CLI based on given arguments.
func (app *App) Run(args []string) error {
	commandPath, flags := splitPathAndFlagValues(app.flagPrefixes, args)
	context := Context{
		App:   app,
		Flags: flags,
	}

	err := executeCommandChain(*(app.baseCommand), commandPath, &context)
	if err != nil {
		return err
	}

	return nil
}

// FlagType Register flag type and it's prefixes.
func (app *App) FlagType(name string, prefixes ...string) {
	for _, prefix := range prefixes {
		app.flagPrefixes[prefix] = name
	}
}

// Command Add command on app.
func (app *App) Command(path string, description string, action CommandAction) *Command {
	return app.baseCommand.Command(path, description, action)
}

// Flag Add flag on app.
func (app *App) Flag(flagType string, name string, description string, action FlagAction) *Flag {
	return app.baseCommand.Flag(flagType, name, description, action)
}

func splitPathAndFlagValues(flagPrefixes map[string]string, args []string) (commandPath []string, flags FlagValues) {
	flags = FlagValues{}
	var currentFlag, flagKey, currentPrefix string
	for _, arg := range args {
		for prefix := range flagPrefixes {
			if strings.HasPrefix(arg, prefix) {
				flagHaveNoValue := currentFlag != "" && currentFlag != arg
				if flagHaveNoValue {
					flags[flagKey] = "true"
				}

				isShorterPrefix := len(prefix) < len(currentPrefix)
				if isShorterPrefix {
					continue
				}

				currentPrefix = prefix
				flagKey = flagPrefixes[prefix] + ":" + arg[len(prefix):]
				currentFlag = arg
			}
		}

		if len(currentFlag) == 0 {
			commandPath = append(commandPath, arg)
		} else if currentFlag != arg {
			flags[flagKey] = arg
			currentFlag = ""
		}
	}

	lastFlagWithNoValue := currentFlag != ""
	if lastFlagWithNoValue {
		flags[flagKey] = "true"
	}

	return
}

func executeCommandChain(parentCommand Command, commandPath []string, context *Context) error {
	for _, command := range parentCommand.Commands {
		argPos := 0

		for _, path := range command.Path {
			isFlagBinder := strings.HasPrefix(path, "{{") && strings.HasSuffix(path, "}}")
			commandPathDontMatch := commandPath[argPos] != path && !isFlagBinder

			if commandPathDontMatch {
				argPos = 0
				break
			}

			if isFlagBinder {
				flagKey := path[2 : len(path)-2]
				context.Flags[flagKey] = commandPath[argPos]
			}

			argPos++
		}

		partialMatch := argPos > 0
		if partialMatch {
			err := executeFlagActions(*command, context)
			if err != nil {
				return err
			}

			if command.Action != nil {
				err = command.Action(context)
				if err != nil {
					return err
				}
			}

			allCommandPathProcessed := argPos == len(commandPath)
			if allCommandPathProcessed {
				return nil
			}

			unprocessedPath := commandPath[argPos:]
			return executeCommandChain(*command, unprocessedPath, context)
		}
	}

	return fmt.Errorf("Command not found")
}

func executeFlagActions(command Command, context *Context) error {
	for rawKey, rawValue := range context.Flags {
		flag, found := command.Flags.getFlag(rawKey)

		if found && flag.Action != nil {
			context.Flags[rawKey] = rawValue
			flag.Action(rawValue, context)
		}

	}
	return nil
}
