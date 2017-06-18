package command

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
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
	log.Println("source and target", source, target)
	if source == "" || target == "" {
		return cli.NewExitError(strings.Join([]string{"source directory path or target json file path is empty\n\nUsage:\n", UsageTextGenarate}, ""), 10)
	}

	list, err := lib.GenerateHashList(source)
	if err != nil {
		return cli.NewExitError(strings.Join([]string{"Failed to genarate hash list. ", err.Error()}, ""), 11)
	}

	data, err := json.MarshalIndent(list, "", "    ")
	if err != nil {
		return cli.NewExitError(strings.Join([]string{"Failed to marshal hash list. ", err.Error()}, ""), 12)
	}

	err = ioutil.WriteFile(target, data, os.ModePerm)
	if err != nil {
		return cli.NewExitError(strings.Join([]string{"Failed to write hash list. ", err.Error()}, ""), 13)
	}
	return nil
}
