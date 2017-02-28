// List command: List all files tored on the network
// Originally work created on 2/4/2017
//

package cmd

import (
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/log"
)

type ListCommand struct{
}

func init() {
	listCmd := ListCommand{}
        config.Parser.AddCommand("list","List files stored in Archit network[Free]", "", &listCmd)
}

func (ec *ListCommand) Execute(args []string) error {
	config.Conf(false)
	log.Console("List command not implemented")
	return nil
}
