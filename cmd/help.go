// Settings command:  Display all configurable settings
// Originally work created on 1/8/2017
//

package cmd

import (
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/log"
)

type HelpCommand struct {
}

var helpCmd HelpCommand

func init() {
	config.Parser.AddCommand("help", "Help using --help", 
		"Really?  You need Help using the Help command?", &helpCmd)
}

func (ec *HelpCommand) Execute(args []string) error {
	config.Conf(false)
	log.Console("'archit help' will display this help\n")
	log.Console("'archit -h' will display the list of all general command line options and commands\n")
	log.Console("'archit <command> -h' will display the general help PLUS any command specific options\n")
	log.Console("Suggest trying a 'archit send -h' as a good example.\n")
	log.Console("Note: -h, -help, and --help are all the same.")
	return nil
}
