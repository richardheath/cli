package cli

// Command Definition of an executable command.
// Commands are matched based on their path.
type Command struct {
	Path      []string // Command path relative to its parent.
	Commands  []Command
	FlagTypes []FlagType
	Usage     string
	Action    CommandAction
}

// CommandAction Method that will envoked when running a command.
type CommandAction func(app App, knownFlags map[string]string, unknownFlags map[string]string) error
