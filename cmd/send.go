// Send command:  Send a file to the network
// Originally work created on 2/4/2017
//

package cmd

import (
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/log"
)

type SendCommand struct{
}

func init() {
	sendCmd := SendCommand{}
        config.Parser.AddCommand("send","Send a file to the Archit network[Fee!]", "", &sendCmd)
}

func (ec *SendCommand) Execute(args []string) error {
	config.Conf(true)
	log.Console("Send command not implemented")
	return nil
}
