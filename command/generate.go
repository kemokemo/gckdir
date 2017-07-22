package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/KemoKemo/gckdir/lib"
	"github.com/urfave/cli"
)

var (
	// UsageGenerate is Usage of generate subcommand for cli
	UsageGenerate = "Generates a json file of the hash list"
)

// CmdGenerate generates a json file of the hash list. This hash list
// includes hash values and recursive directory structure information.
func CmdGenerate(c *cli.Context) error {
	help := fmt.Sprintf("Please see '%s %s --help'.", c.App.Name, c.Command.FullName())
	source := c.Args().Get(0)
	target := c.Args().Get(1)
	if source == "" || target == "" {
		return cli.NewExitError(
			fmt.Sprintf("Source path or target path is empty. %s", help),
			ExitCodeInvalidArguments)
	}
	source = filepath.Clean(source)
	target = filepath.Clean(target)

	list, err := lib.GetHashList(source)
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("Failed to genarate hash list. %v\n%s", err, help),
			ExitCodeFunctionError)
	}

	data, err := json.MarshalIndent(list, "", "    ")
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("Failed to marshal hash list. %v\n%s", err, help),
			ExitCodeFunctionError)
	}

	err = ioutil.WriteFile(target, data, os.ModePerm)
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("Failed to write hash list. %v\n%s", err, help),
			ExitCodeIOError)
	}
	return nil
}
