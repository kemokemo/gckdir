package main

import (
	"fmt"
	"os"

	"github.com/kemokemo/gckdir/command"
	"github.com/urfave/cli"
)

func main() {
	os.Exit(run(os.Args))
}

func run(args []string) int {
	app := cli.NewApp()
	app.Name = Name
	app.Version = fmt.Sprintf("%s.%s", Version, Revision)
	app.Author = "kemokemo"
	app.Email = "t2wonderland@gmail.com"
	app.Usage = "generate a hash list of a correct directory and verify the target directory's structure and each hash value of files."

	app.Flags = GlobalFlags
	app.Commands = Commands
	app.CommandNotFound = CommandNotFound

	err := app.Run(args)
	if err != nil {
		fmt.Println("failed to run: ", err)
		return command.ExitCodeFailed
	}
	return command.ExitCodeOK
}
