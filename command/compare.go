package command

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/kemokemo/gckdir/lib"
	"github.com/urfave/cli"
)

var (
	// UsageCompare is Usage of compare subcommand for cli
	UsageCompare = "Compares directory information"

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
		return cli.NewExitError(fmt.Sprintf("source path or target path is empty\n\nUsage:\n%s", UsageTextCompare), ExitCodeInvalidArguments)
	}
	source = filepath.Clean(source)
	target = filepath.Clean(target)

	sourceList, err := lib.GetHashList(source)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("failed to get the hash list of '%s'.\n\nUsage:\n%s", source, UsageTextCompare), ExitCodeFunctionError)
	}

	targetList, err := lib.GetHashList(target)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("failed to get the hash list of '%s'.\n\nUsage:\n%s", target, UsageTextCompare), ExitCodeFunctionError)
	}

	log.Println("Source:", source)
	log.Println("Target:", target)
	result := lib.CompareHashList(sourceList, targetList)
	if result.CompareResult {
		log.Println("The comparison was successful.")
	} else {
		log.Println("The comparison failed.")
	}
	return nil
}
