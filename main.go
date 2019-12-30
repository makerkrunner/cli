// +build go1.13

package main

import (
	command_parser "code.cloudfoundry.org/cli/parser"
	"code.cloudfoundry.org/cli/plugin_parser"
	"code.cloudfoundry.org/cli/util/panichandler"
	"os"
)

func main() {
	var exitCode int
	defer panichandler.HandlePanic()
	plugin, commandIsPlugin := plugin_parser.IsPluginCommand(os.Args[1])

	if commandIsPlugin == true {
		exitCode = plugin_parser.RunPlugin(plugin)
	} else {
		exitCode = command_parser.ParseCommandFromArgs(os.Args)
	}
	if exitCode != 0 {
		os.Exit(exitCode)
	}
}
