package cli

import "strings"

// App CLI Application
type App struct {
	Name                    string
	Version                 string
	FlagPrefixes            []FlagPrefix
	FlagTypes               []FlagType
	FlagValidationFormatter FlagValidationFormatter
	Commands                []Command
	CommandArgs             []string
	FlagArgs                []string
}

// NewApp Initialize App struct.
func NewApp(name string, version string) App {
	return App{
		Name:                    name,
		Version:                 version,
		FlagPrefixes:            []FlagPrefix{},
		FlagTypes:               []FlagType{},
		FlagValidationFormatter: DefaultFlagValidationFormatter,
		Commands:                []Command{},
		CommandArgs:             []string{},
		FlagArgs:                []string{},
	}
}

// Run Execute command based on given arguments.
func (app *App) Run(args []string) error {
	var err error
	app.ProcessArgs(args)

	match, err := app.GetMatchingCommand()
	if err != nil {
		return err
	}

	flags, err := app.ProcessFlags(match)
	if err != nil {
		return err
	}

	insertFlagsDefaultValues(flags, match)
	err = match.Command.Action(flags)
	if err != nil {
		return err
	}

	return nil
}

// ProcessArgs Categorize command/flag arguments.
func (app *App) ProcessArgs(args []string) {
	flagPrefixes := getFlagPrefixes(app)

	var currentFlag string
	for _, arg := range args {
		for prefix := range flagPrefixes {
			if strings.HasPrefix(arg, prefix) {
				currentFlag = arg
				break
			}
		}

		if len(currentFlag) == 0 {
			app.CommandArgs = append(app.CommandArgs, arg)
		} else {
			app.FlagArgs = append(app.FlagArgs, arg)

			if currentFlag != arg {
				currentFlag = ""
			}
		}
	}
}

// GetMatchingCommand Get matching command using processed arguments.
// Must call ProcessArgs first.
func (app *App) GetMatchingCommand() (CommandMatchInfo, error) {
	matchInfo := CommandMatchInfo{
		unprocessedArgs: app.CommandArgs[:],
		BindedFlags:     []string{},
		FlagTypes:       []FlagType{},
	}

	err := getMatchingCommand(app.Commands, &matchInfo)
	return matchInfo, err
}

func getFlagPrefixes(app *App) map[string]string {
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
