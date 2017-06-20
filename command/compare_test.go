package command

import (
	"flag"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/urfave/cli"
)

func TestCmdCompare(t *testing.T) {
	flagPass := flag.FlagSet{}
	err := flagPass.Parse(strings.Split(
		fmt.Sprintf("%s %s", filepath.Join(dir, "TestData/dir1.json"), filepath.Join(dir, "TestData/Dir1")), " "))
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
			if err = CmdCompare(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("CmdCompare() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
