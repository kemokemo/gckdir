package command

import "github.com/urfave/cli"

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
	// Write your code here
	return nil
}
