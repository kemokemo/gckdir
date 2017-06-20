package main

import (
	"fmt"
	"os"

	"github.com/KemoKemo/gckdir/command"
	"github.com/urfave/cli"
)

var (
	// GlobalFlags are global flag values
	GlobalFlags = []cli.Flag{}

	// Commands are sub-commands of this app
	Commands = []cli.Command{
		{
			Name:      "generate",
			Aliases:   []string{"gen"},
			Usage:     command.UsageGenerate,
			UsageText: command.UsageTextGenarate,
			Action:    command.CmdGenerate,
			Flags:     []cli.Flag{},
			OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
				fmt.Fprintf(c.App.Writer, "%v\nPlease check the usage below.\n\n%s\n", err, command.UsageTextGenarate)
				return err
			},
		},
		{
			Name:    "compare",
			Aliases: []string{"com"},
			Usage:   "",
			Action:  command.CmdCompare,
			Flags:   []cli.Flag{},
			OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
				fmt.Fprintf(c.App.Writer, "%v\nPlease check the usage below.\n\n%s\n", err, command.UsageTextCompare)
				return err
			},
		},
	}
)

// CommandNotFound will be executed when the user inputed sub-command is invalid.
func CommandNotFound(c *cli.Context, subcommand string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.\n", c.App.Name, subcommand, c.App.Name, c.App.Name)
	os.Exit(command.ExitCodeCommandNotFound)
}
