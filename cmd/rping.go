// Status command:  Display Network and Wallet status
// Originally work created on 1/8/2017
//

package cmd

import (
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/log"
	"github.com/valyala/gorpc"
	"fmt"
	"net"
)

type RPingCommand struct {
}

func init() {
	pingCmd := RPingCommand{}
	config.Parser.AddCommand("rping", "Pings a remote farming service[Free]", "", &pingCmd)
}

func (ec *RPingCommand) Execute(args []string) error {
	config.Conf(false)

	/// Insert RPC code to query the farmer

	serverIP := "127.0.0.1:1958"
	fmt.Printf("Please enter a remote Host IP address in the form 1.2.3.4:1958 ")
	_, err := fmt.Scanf("%s", &serverIP)
	if err != nil {
		log.Critical("I'm sorry Dave, I'm afraid I can't do that")
	}
	host,port,err := net.SplitHostPort(serverIP)
	if err != nil {
		log.Critical("Invalid input",serverIP)
	}
	log.Console("Attempting to connect to server",host,"on port",port)
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
