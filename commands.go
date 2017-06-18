package main

import (
	"fmt"
	"os"

	"github.com/KemoKemo/gckdir/command"
	"github.com/urfave/cli"
)

var (
	GlobalFlags = []cli.Flag{}
	Commands    = []cli.Command{
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

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}
