package command

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"

	"github.com/kemokemo/gckdir/lib"
	"github.com/urfave/cli"
)

var (
	// UsageCompare is Usage of compare subcommand for cli
	UsageCompare = "Generates a json file of the hash list"

	// UsageTextCompare is UsageText of compare subcommand for cli
	UsageTextCompare = `Compares directory information below cases.
 Case 1. a json file of hash list with target directory
 Case 2. source directory with target directory

 Usage: gckdir compare [source directory path or a json file path of hash list] [target directory]
   ex1) gckdir compare path/to/source path/to/target
	 ex2) gckdir compare hash.json path/to/target`
)

// CmdCompare comares directory information below cases.
//  Case 1. a json file of hash list with target directory
//  Case 2. source directory with target directory
func CmdCompare(c *cli.Context) error {
	source := c.Args().Get(0)
	target := c.Args().Get(1)
	if source == "" || target == "" {
		return cli.NewExitError(strings.Join([]string{"source path or target path is empty\n\nUsage:\n", UsageTextCompare}, ""), 20)
	}

	// TODO: lib.ReadHashList(source)
	data, err := ioutil.ReadFile(source)
	if err != nil {
		return cli.NewExitError(strings.Join([]string{"Failed to read a json file.\n\nUsage:\n", UsageTextCompare}, ""), 21)
	}

	sourceList := lib.HashList{}
	err = json.Unmarshal(data, &sourceList)
	if err != nil {
		return cli.NewExitError(strings.Join([]string{"Failed to unmarshal.\n\nUsage:\n", UsageTextCompare}, ""), 21)
	}

	targetList, err := lib.GenerateHashList(target)
	if err != nil {
		return cli.NewExitError(strings.Join([]string{"Failed to genarate hash info.\n\nUsage:\n", UsageTextCompare}, ""), 21)
	}

	result := lib.CompareHashList(sourceList, targetList)
	if result.CompareResult {
		log.Printf("Successfully compared.\n source:%s\n target:%s", source, target)
	} else {
		log.Printf("Failed to compare.\n source:%s\n target:%s", source, target)
	}
	return nil
}
