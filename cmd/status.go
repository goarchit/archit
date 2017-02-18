// Status command:  Display Network and Wallet status
// Originally work created on 1/8/2017
//

package cmd

import (
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/farmer"
	"github.com/valyala/gorpc"
	"errors"
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
	if !farmer.FarmerStarted {
		return errors.New("Farming has not been started; Please issue a 'archit farm' command")
	}
	/// Insert RPC code to query the farmer

	port := config.Archit.PortBase + 1
	serverIP := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
	c := gorpc.NewTCPClient(serverIP)
	c.Start()
	defer c.Stop()
	c.Call("Status")
	return nil
}
