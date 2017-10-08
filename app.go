package cli

import (
	"fmt"
	"strings"
)

// App Represents CLI application.
type App struct {
	Name         string
	Version      string
	baseCommand  *Command
	flagPrefixes map[string]string
}

// Context Context of current CLI execution.
type Context struct {
	App          *App
	Flags        FlagValues
	Data         interface{}
	commandChain []*Command
}

// NewApp Returns a new CLI application.
func NewApp(name, version string) App {
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
		App:          app,
		Flags:        flags,
		commandChain: []*Command{app.baseCommand},
	}

	err := detectCommandChain(*(app.baseCommand), commandPath, &context)
	if err != nil {
		return err
	}

	return executeCommandChain(&context)
}

// FlagType Register flag type and it's prefixes.
func (app *App) FlagType(name string, prefixes ...string) {
	for _, prefix := range prefixes {
		app.flagPrefixes[prefix] = name
	}
}

// Command Add command on app.
func (app *App) Command(path string, action CommandAction) *Command {
	return app.baseCommand.Command(path, action)
}

// Flag Add flag on app.
func (app *App) Flag(flagType, name, defaultValue string, action FlagAction) *Flag {
	return app.baseCommand.Flag(flagType, name, defaultValue, action)
}

func splitPathAndFlagValues(flagPrefixes map[string]string, args []string) (commandPath []string, flags FlagValues) {
	flags = FlagValues{}
	var currentFlag, flagKey, currentPrefix string
	for _, arg := range args {
		for prefix, flagType := range flagPrefixes {
			if strings.HasPrefix(arg, prefix) {
				flagHaveNoValue := currentFlag != "" && currentFlag != arg
				if flagHaveNoValue {
					flags[flagKey] = ""
				}

				isShorterThanLastPrefix := len(prefix) < len(currentPrefix)
				if isShorterThanLastPrefix {
					continue
				}

				currentPrefix = prefix
				flagKey = flagType + ":" + arg[len(prefix):]
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

func detectCommandChain(parentCommand Command, unprocessedArgs []string, context *Context) error {
	for _, command := range parentCommand.Commands {
		argPos := 0
		for _, path := range command.Path {
			isFlagBinder := strings.HasPrefix(path, "{{") && strings.HasSuffix(path, "}}")
			commandPathDontMatch := unprocessedArgs[argPos] != path && !isFlagBinder

			if commandPathDontMatch {
				argPos = 0
				break
			}

			if isFlagBinder {
				flagKey := path[2 : len(path)-2]
				context.Flags[flagKey] = unprocessedArgs[argPos]
			}

			argPos++
		}

		if partialMatch := argPos > 0; partialMatch {
			context.commandChain = append(context.commandChain, command)

			if allArgsProcessed := argPos == len(unprocessedArgs); allArgsProcessed {
				return nil
			}

			unprocessedPath := unprocessedArgs[argPos:]
			return detectCommandChain(*command, unprocessedPath, context)
		}
	}

	return fmt.Errorf("Command not found")
}

func executeCommandChain(context *Context) error {
	for _, command := range context.commandChain {
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
	}

	return nil
}

func executeFlagActions(command Command, context *Context) error {
	for _, flag := range command.Flags {
		flagValue := ""
		if value, matchFlagName := context.Flags[flag.nameKey()]; matchFlagName {
			flagValue = value
		} else if value, matchFlagShorthand := context.Flags[flag.shorthandKey()]; matchFlagShorthand {
			flagValue = value
		} else if flag.Default != "" {
			flagValue = flag.Default
		} else {
			continue
		}

		if flagValue == "" && flag.Default != "" {
			flagValue = flag.Default
		}

		context.Flags[flag.nameKey()] = flagValue
		if flag.Action != nil {
			err := flag.Action(flagValue, context)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
