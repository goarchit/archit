// Create command:  Interactively create the configuration file
// Originally work created on 2/4/2017
//

package cmd

import (
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/log"
)

type CreateCommand struct {
}

func init() {
	createCmd := CreateCommand{}
	config.Parser.AddCommand("create", "Interactively create/edit your configuration file", "", &createCmd)
}

func (ec *CreateCommand) Execute(args []string) error {
	config.Conf(false)
	log.Console("Create command not implemented")
	return nil
}
