package cli

type Command struct {
	Path       string // Command path relative to its parent.
	Commands   []Command
	FlagGroups []FlagGroup
	FlagTypes  []FlagType
	Action     CommandAction
}

type CommandAction func(knownFlags map[string]string, unknownFlags map[string]string) error
