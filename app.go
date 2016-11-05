package cli

import (
	"fmt"
	"os"
	"strings"
)

type App struct {
	Name       string
	Version    string
	FlagGroups []FlagGroup
	FlagTypes  []FlagType
	Commands   []Command
}

func (app App) Run(args []string) error {
	commandPath, err := getCommandPath(app, args)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(commandPath)
	// getCommandPath from args
	// getMatching command
	// if no match return error

	// get normalized flags
	// process flags using app and command flag options
	// if error return error

	// execute command action
	// if error return error
	return nil
}

func getCommandPath(app App, args []string) ([]string, error) {
	prefixes := getFlagPrefixes(app)
	command := make([]string, 0, 5)

	for _, arg := range args {
		flag := false
		for _, prefix := range prefixes {
			if strings.HasPrefix(arg, prefix) {
				flag = true
			}
		}

		if !flag {
			command = append(command, arg)
		}
	}
	return command, nil
}

func getFlagPrefixes(app App) []string {
	prefixes := make([]string, 0, 2)
	for _, group := range app.FlagGroups {
		if group.Prefix != "" {
			prefixes = append(prefixes, group.Prefix)
		}

		if group.ShorthandPrefix != "" {
			prefixes = append(prefixes, group.ShorthandPrefix)
		}
	}

	return prefixes
}

func contains(arr []string, value string) bool {
	for _, item := range arr {
		if item == value {
			return true
		}
	}
	return false
}
