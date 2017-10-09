# Go CLI Library

Simple but powerful Go CLI Library. The design favors simple patterns that 
you can expand to create simple to complex CLI that fits your needs.

```go
func main() {
  app := cli.NewApp("mycli", "1.0.0")
  app.FlagType("option", "--", "-")

  app.Command("greet", func(ctx *cli.Context) error {
    fmt.Println("hello")
    return nil
  })
}
```

## Features

### Custom Flag Types

Custom flag types that your CLI supports. This gives capability
to use different flag groups for different purpose.
```go
/*
app.FlagType(name, keys...)
* name - Flag type name. This is used to retreive flag values.
* keys - Flag prefixes for given flag.
*/
app.FlagType("option", "--", "-")
app.FlagType("setting", "#")
```

### Rest like command path definition

Command path defines how to get to the target command.

```go
// args to invoke: greet
app.Command("greet", genericGreet)
// args to invoke: email group
app.Command("email group", emailGroup)
// args to invoke: email user richard
// this is using flag binding feature
app.Command("email user {{option:user}}", emailUser)
```

### Nested Commands

Commands can be nested. When args invokes child commands it
will create a chain where parent commands will get invoked
first before invoking child command. This way you can specify
function that is common to all child commands. If you don't
want to run anything just set the command action to `nil`.

On sample below args `email user manager` will invoke
`emailUser` and `emailUserManager`.

```go
email := app.Command("email", nil)
// args to invoke: email group
email.Command("group", emailGroup)
// args to invoke: email user
emailUser := email.Command("user", emailUser)
// args to invoke: email user manager
emailUser.Command("manager", emailUserManager)
```

### Nested Flags

Given the command chain. All parent flags are inherited
by child flag. It is also executed based on command chain
order. On sample below option email validation is defined
only on parent but will always get executed on sub commands.

```go
email := app.Command("email", nil)
email.flag("option", "email", "", func(value string, ctx *cli.Context) error {
  err := validEmail(value)
  if err != nil {
    return err
  }

  return nil
})

email.Command("group", emailGroup)
email.Command("user", emailUser)
```

### Flag binding

This feature allows values to be binded to flags
from command path. This makes it easy to create
pretty CLI syntax like:

> print --message hello

This can be written as:

> print hello

```go
// args to invoke: email user richard
// where richard will be binded to flag user
app.Command("email user {{option:user}}", emailUser)
```

## Actions

Flag and command definition can have actions. These
functions are invoked if they are part of the 
command chain. Supply `nil` functions to ignore.

When action returns an error the chain will stop
executing and return the error.

```go
// Define user flag with validation.
app.flag("option", "user", "", func(value string, ctx *cli.Context) error {
  err := validUser(value)
  if err != nil {
    return err
  }

  return nil
})

// Send email to validated user.
app.Command("email user {{option:user}}", func(ctx *cli.Context) error {
  user, _ := ctx.Flags.GetValue("option", "user")
  err := sendImportandEmailToUser(user)
  return err
})
```

## Sample App

``` go
package main

import (
  "fmt"
  "os"
  "github.com/richardheath/cli"
)

func main() {
  app := cli.NewApp("app", "0.1.0")
  app.FlagType("option", "--", "-")

  // Global flags
  app.Flag("option", "log l", "console", initLogging)

  // Command definition
  greet := app.Command("greet {{option:user}}", func(ctx *cli.Context) error {
    user, _ := ctx.Flags.GetValue("option", "user")
    fmt.Printf("hello %s", user)
    return nil
  })

  // Flag only available on greet.
  greet.Flag("option", "user u", nil)

  // Run the application based on args
  err := app.Run(os.Args[1:])
  if err != nil {
    fmt.Println(err.Error())
  }
}

```