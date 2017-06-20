package command

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/kemokemo/gckdir/lib"
	"github.com/urfave/cli"
)

var (
	// UsageGenerate is Usage of generate subcommand for cli
	UsageGenerate = "Generates a json file of the hash list"

	// UsageTextGenarate is UsageText of generate subcommand for cli
	UsageTextGenarate = `Generates a json file of the hash list.
This hash list includes hash values and recursive directory structure information.
 Usage: gckdir generate [directory path] [json file path or name]
    ex) gckdir generate path/to/source source.json`
)

// CmdGenerate generates a json file of the hash list. This hash list
// includes hash values and recursive directory structure information.
func CmdGenerate(c *cli.Context) error {
	source := c.Args().Get(0)
	target := c.Args().Get(1)
	if source == "" || target == "" {
		return cli.NewExitError(strings.Join([]string{"source path or target path is empty\n\nUsage:\n", UsageTextGenarate}, ""), ExitCodeInvalidArguments)
	}
	source = filepath.Clean(source)
	target = filepath.Clean(target)

	list, err := lib.GetHashList(source)
	if err != nil {
		return cli.NewExitError(strings.Join([]string{"Failed to genarate hash list. ", err.Error()}, ""), ExitCodeFunctionError)
	}

	data, err := json.MarshalIndent(list, "", "    ")
	if err != nil {
		return cli.NewExitError(strings.Join([]string{"Failed to marshal hash list. ", err.Error()}, ""), ExitCodeFunctionError)
	}

	err = ioutil.WriteFile(target, data, os.ModePerm)
	if err != nil {
		return cli.NewExitError(strings.Join([]string{"Failed to write hash list. ", err.Error()}, ""), ExitCodeIOError)
	}
	return nil
}
