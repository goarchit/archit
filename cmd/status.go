// Status command:  Display Network and Wallet status
// Originally work created on 1/8/2017
//

package cmd

import (
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/log"
	"github.com/valyala/gorpc"
	"net"	
	"strconv"
)

type StatusCommand struct{
}

func init() {
	statusCmd := StatusCommand{}
        config.Parser.AddCommand("status","Shows the status of the ArchIt farming & wallet servers[Free]", "", &statusCmd)
}

func (ec *StatusCommand) Execute(args []string) error {
	config.Conf(false)

	/// Insert RPC code to query the farmer

	port := config.Archit.PortBase + 1
	serverIP := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
	c := gorpc.NewTCPClient(serverIP)
	c.Start()
	defer c.Stop()

	d := gorpc.NewDispatcher()
	d.AddFunc("Status", func() {})
	dc := d.NewFuncClient(c)
	response, err := dc.Call("Status", nil)
	if err != nil {
                log.Error("Status failed: ",err)
        }
	log.Console(response)
	return nil
}
