package cli

// FlagGroup Definition for flag groups. This can be used to specify different types of flag groups that the app supports.
type FlagGroup struct {
	Prefix          string
	ShorthandPrefix string
	Group           string
}

type FlagGroupValues struct {
	Known   map[string]string
	Unknown map[string]string
}

type FlagType struct {
	key          string
	defaultValue string
	validators   []FlagValidator
}

type FlagValidator func(key string, value string) error
