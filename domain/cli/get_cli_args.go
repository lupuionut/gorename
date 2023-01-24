package cli

import (
    "fmt"
    "strings"
    "errors"
)

type Instance struct {
    Commands map[string]string
    Error error
}

type Severity int
const (
    Low = 0
    Fatal = 1
)

type CommandParseError struct {
    Err error
    Level Severity
}

func (e *CommandParseError) Error() string {
    return fmt.Sprintf("%s", e.Err)
}

func (cli *Instance) ParseArgs(args []string) {
    if len(args) < 2 {
        cli.Error = &CommandParseError {
            Err: errors.New("Not enogh command line arguments provided.\n" + cli.Help()),
            Level: Fatal,
        }
        return
    }

    var errorText string
    args = args[1:]

    for i := range(args) {
        exp := args[i]
        if string(exp[0]) != "-" {
            errorText += fmt.Sprintf("Could not parse the following argument '%s'. Make sure you use the format -key=value \n", exp)
            continue
        }

        parts := strings.Split(exp, "=")
        if len(parts) != 2 {
            errorText += fmt.Sprintf("Invalid format for argument declaration for '%s'. You must specify each argument as -key=value \n", exp)
            continue
        }
        key := string(parts[0][1:])
        value := string(parts[1])
        cli.Commands[key] = value
    }

    if errorText != "" {
        cli.Error = &CommandParseError{
            Err: fmt.Errorf(errorText),
            Level: Low,
        }
    }
}

func (cli *Instance) Help() string {
    helpText := "Usage: gorename COMMANDS \n\n COMMANDS:\n"
    for v := range(cli.Commands) {
        var choices string
        if v == "path" {
            choices = "You must specify the full path to the folder that contains the files to rename."
        }
        helpText += fmt.Sprintf("`-%s= `. " + choices + "\n", v)
    }
    return helpText
}
