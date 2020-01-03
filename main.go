// +build go1.13

package main

import (
	"os"

	"code.cloudfoundry.org/cli/cf/cmd"
	command_parser "code.cloudfoundry.org/cli/parser"
	"code.cloudfoundry.org/cli/plugin_parser"
	"code.cloudfoundry.org/cli/util/panichandler"
)

const unknownCommandCode = -666

func main() {
	var exitCode int
	defer panichandler.HandlePanic()

	exitCode = command_parser.ParseCommandFromArgs(os.Args)
	if exitCode == unknownCommandCode {
		plugin, commandIsPlugin := plugin_parser.IsPluginCommand(os.Args)

		if commandIsPlugin == true {
			exitCode = plugin_parser.RunPlugin(plugin)
		} else {
			cmd.Main(os.Getenv("CF_TRACE"), os.Args)
		}
	}

	if exitCode != 0 {
		os.Exit(exitCode)
	}
}
