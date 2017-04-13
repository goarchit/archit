// Get command:  fetch a file from the network
// Originally work created on 1/8/2017
//

package cmd

import (
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/log"
	"fmt"
	"os"
)

type GetCommand struct {
}

func init() {
	getCmd := GetCommand{}
	_,err := config.Parser.AddCommand("get", "Retrieve a file from the Archit network[Fee!]", "", &getCmd)
        if err != nil {
                fmt.Println("Internal error parsing Get command:",err)
                os.Exit(1)
        }
}

func (ec *GetCommand) Execute(args []string) error {
	config.Conf(true)
	log.Console("Get command not implemented")
	return nil
}
