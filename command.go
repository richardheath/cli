package cli

type Command struct {
	Path      []string // Command path relative to its parent.
	Commands  []Command
	FlagTypes []FlagType
	Action    CommandAction
}

type CommandAction func(knownFlags map[string]string, unknownFlags map[string]string) error
