// Ping  command:  Validate a service is running
// Originally work created on 1/8/2017
//

package cmd

import (
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/log"
	"github.com/valyala/gorpc"
	"fmt"
	"net"
	"strconv"
)

type PingCommand struct {
	Internal bool `short:"i" long:"Internal" description:"Ping the internal server"`
	External bool `short:"e" long:"External" description:"Ping the public facing server"`
	TCP bool `short:"t" long:"TCP" description:"Ping server via TCP address"`
}

var pingCmd PingCommand
var port int
var serverIP string

func init() {
	config.Parser.AddCommand("ping", "Ping a service, expect a Pong", "", &pingCmd)
}

func (ec *PingCommand) Execute(args []string) error {
	config.Conf(false)

	// Set default
	if !pingCmd.Internal && !pingCmd.External && !pingCmd.TCP {
		pingCmd.Internal = true
	}

	/// Insert RPC code to query the farmer

	if pingCmd.Internal {
		port = config.Archit.PortBase + 1
		serverIP = net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
	} else if pingCmd.External {
		port = config.Archit.PortBase
        	serverIP = net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
	} else if pingCmd.TCP {
        	fmt.Printf("Please enter a remote Host IP address in the form 1.2.3.4:1958 ")
        	_, err := fmt.Scanf("%s", &serverIP)
        	if err != nil {
               		log.Critical("I'm sorry Dave, I'm afraid I can't do that")
        	}
        	host,sport,err := net.SplitHostPort(serverIP)
        	if err != nil {
               		log.Critical("Invalid input",serverIP)
        	}
        	log.Console("Attempting to connect to server",host,"on port",sport)
	}
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
