# Go CLI Library

Dynamic Go CLI Library. Commands can have their own set of commands. This allows 
Sub-commands inherits flag tpyes of it's parents. Flag types defined on app will
be available to all commands.

The cli library also supports flag binding. This allows values to be binded to
flags based semantics like.

> print --message hello

This can be written as:

> print hello

## Sample Binding Command


``` go
package main

import (
    "fmt"
    "os"
    "github.com/richardheath/cli"
)

func main() {
    app := cli.NewApp("app", "0.1.0")
    app.FlagPrefixes = []cli.FlagPrefix{
      cli.FlagPrefix{
        Key:         "--",
        Shorthand:   "-",
        Description: "options",
      },
    }

    app.FlagTypes = []cli.FlagType{
      cli.FlagType{
        Key:        "message",
        Shorthand:  "m",
        Prefix:     "--",
      },
    }

    app.Commands = []cli.Command{
      cli.Command{
        Path:  []string{"print", "{{--message}}"},
        Usage: "print {message}",
        Action: func(flags cli.ProcessedFlags) error {
          fmt.Println(flags.known["--message"])
          return nil
        },
      },
    },
  }

  app.Run(os.Args[1:])
}
```