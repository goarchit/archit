// Get command:  fetch a file from the network
// Originally work created on 1/8/2017
//

package cmd

import (
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/log"
)

type GetCommand struct{
}

func init() {
	getCmd := GetCommand{}
        config.Parser.AddCommand("get","Retrieve a file from the Archit network[Fee!]", "", &getCmd)
}

func (ec *GetCommand) Execute(args []string) error {
	config.Conf(true)
	log.Console("Get command not implemented")
	return nil
}
