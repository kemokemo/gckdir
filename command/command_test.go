package command

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/urfave/cli"
)

var (
	flags    = []cli.Flag{}
	commands = []cli.Command{
		{
			Name:   "generate",
			Action: CmdGenerate,
			Flags:  []cli.Flag{},
			OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
				fmt.Fprintf(c.App.Writer, "%v\nPlease check the usage below.\n\n%s\n", err, UsageTextGenarate)
				return err
			},
		},
		{
			Name:   "compare",
			Action: CmdCompare,
			Flags:  []cli.Flag{},
			OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
				fmt.Fprintf(c.App.Writer, "%v\nPlease check the usage below.\n\n%s\n", err, UsageTextCompare)
				return err
			},
		},
	}
	app = cli.NewApp()
	dir = ""
)

func TestMain(t *testing.M) {
	setup()
	exitCode := t.Run()
	os.Exit(exitCode)
}

func setup() {
	app.Flags = flags
	app.Commands = commands

	var err error
	dir, err = os.Getwd()
	if err != nil {
		log.Println("Failed to get current directory path. ", err)
	}
}
