package command

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
)

func TestCmdGenerate(t *testing.T) {
	app := cli.NewApp()
	app.Flags = flags
	app.Commands = commands

	dir, err := os.Getwd()
	if err != nil {
		t.Errorf("Failed to get current directory path. %v", err)
	}

	flagPass := flag.FlagSet{}
	err = flagPass.Parse(strings.Split(fmt.Sprintf("%s dir1.json", filepath.Join(dir, "TestData/Dir1")), " "))
	if err != nil {
		t.Errorf("Failed to create args. %v", err)
	}
	flagFail1 := flag.FlagSet{}
	err = flagFail1.Parse(strings.Split(fmt.Sprintf("%s", filepath.Join(dir, "TestData/Dir1")), " "))
	if err != nil {
		t.Errorf("Failed to create args. %v", err)
	}
	flagFail2 := flag.FlagSet{}
	err = flagFail2.Parse(strings.Split("", ""))
	if err != nil {
		t.Errorf("Failed to create args. %v", err)
	}

	type args struct {
		c *cli.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "TestPass-1", args: args{c: cli.NewContext(app, &flagPass, nil)}, wantErr: false},
		{name: "TestFail-1", args: args{c: cli.NewContext(app, &flagFail1, nil)}, wantErr: true},
		{name: "TestFail-2", args: args{c: cli.NewContext(app, &flagFail2, nil)}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err = CmdGenerate(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("CmdGenerate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	err = os.Remove("dir1.json")
	if err != nil {
		t.Errorf("Failed to remove dir1.json file. %v", err)
	}
}
