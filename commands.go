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
			Action:    command.CmdGenerate,
			Flags:     []cli.Flag{},
			ArgsUsage: "[source] [target]\n\t\tsource: a directory path\n\t\ttarget: a json file path",
			OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
				fmt.Fprintf(c.App.Writer, "A ussage error occurred. Please see '%s %s --help'.\n", c.App.Name, c.Command.FullName())
				return err
			},
		},
		{
			Name:    "verify",
			Aliases: []string{"ver"},
			Usage:   command.UsageVerify,
			Action:  command.CmdVerify,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "report, r",
					Usage: "Create a result report in html format.",
				},
				cli.StringFlag{
					Name:  "output, o",
					Usage: "Create a report file with the name specified by this flag. Ignored if --report or --open flag is not specified.",
				},
				cli.BoolFlag{
					Name:  "open, p",
					Usage: "Create and open a result report with the default browser. This option includes the 'report' option.",
				},
				cli.BoolFlag{
					Name:  "no-hv, nh",
					Usage: "Verify only the file and directory structure. Ignore the each hash value of files.",
				},
				cli.BoolFlag{
					Name:  "no-uv, nu",
					Usage: "Ignore the unnecessary files.",
				},
			},
			ArgsUsage: "[source] [target]\n\t\tsource: a json file path or a directory path\n\t\ttarget: a directory path",
			OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
				fmt.Fprintf(c.App.Writer, "A ussage error occurred. Please see '%s %s --help'.\n", c.App.Name, c.Command.FullName())
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
