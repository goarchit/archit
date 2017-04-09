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

type EPingCommand struct {
}

func init() {
	pingCmd := EPingCommand{}
	config.Parser.AddCommand("eping", "Pings the external farming service[Free]", "", &pingCmd)
}

func (ec *EPingCommand) Execute(args []string) error {
	config.Conf(false)

	/// Insert RPC code to query the farmer

	port := config.Archit.PortBase 
	serverIP := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
	c := gorpc.NewTCPClient(serverIP)
	c.Start()
	defer c.Stop()

	d := gorpc.NewDispatcher()
	d.AddFunc("Ping", func() {})
	dc := d.NewFuncClient(c)
	resp, err := dc.Call("Ping", nil)
	if err != nil {
		log.Error("Ping failed: ", err)
	}
	log.Console("Farmer process responded: ", resp)
	return nil
}
