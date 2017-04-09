// Audit command - Performs various Archit Network audits
// Originally work created on 2/4/2017
//

package cmd

import (
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/log"
	"github.com/goarchit/archit/util"
	"net"
)

type TestCommand struct {
}

func init() {
	testCmd := TestCommand{}
	config.Parser.AddCommand("test", "Go functionality testing", "", &testCmd)
}

func (ec *TestCommand) Execute(args []string) error {
	config.Conf(false)
	rifs := util.RoutedInterface("ip", net.FlagUp|net.FlagBroadcast)
	if rifs != nil {
		log.Console("Routed interface is ", rifs.HardwareAddr.String())
	}

	log.Console("Testing complete")

	return nil
}
