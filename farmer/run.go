package farmer

import (
	"github.com/briandowns/spinner"
	"github.com/goarchit/archit/log"
	"github.com/goarchit/archit/util"
	"github.com/valyala/gorpc"
	"net"
	"strconv"
	"time"
)

var dbCmd = make(chan string)
var walletCmd = make(chan string)
var intRPC = gorpc.NewDispatcher()
var extRPC = gorpc.NewDispatcher()
var intCmd *gorpc.Server
var extCmd *gorpc.Server

func Run(c chan bool) {

	// Go look up seeds and set util.IAMASeed
	util.Dnsseed()

	// Start common services
	if !util.IAmASeed && !util.SeedMode {
		go Wallet(walletCmd,true)
	}
	go DB(dbCmd)

	// External RPC service first
	// Start by registering fucntions and types
	extRPC.AddFunc("Ping", func() string { return "ePong!" })
	extRPC.AddFunc("PeerAdd", func(pi *PeerInfo) string { return PeerAdd(pi) })
	// Only seed nodes can be queried for peer info from the outside world
	if util.IAmASeed {
		extRPC.AddFunc("PeerListAll", func() string { return PeerListAll() })
	}

	// Then launch the server
	serverIP := ":" + strconv.Itoa(util.PortBase) // Listen on all interfaces
	log.Info("Farmer External RPC Server using server address", serverIP)
	extCmd = gorpc.NewTCPServer(serverIP, extRPC.NewHandlerFunc())
	err := extCmd.Start()
	if err != nil {
		log.Critical("Farmer External RPC service failed to start: ", err)
	}

	defer extCmd.Stop()

	// Internal RCP service next
	// Start by registering functions and types
	intRPC.AddFunc("Ping", func() string { return "iPong!" })
	intRPC.AddFunc("Status", func() string { return Status() })
	// intRPC.AddFunc("PeerAdd", func(pi *PeerInfo) error {return PeerAdd(wa,p)})
	intRPC.AddFunc("PeerDelete", func(p *util.Peer) error { return PeerDelete(p) })
	intRPC.AddFunc("PeerListAll", func() string { return PeerListAll() })
	intRPC.AddFunc("PeerDBSync", func() error { return FlushPeerMap() })

	// Then launch the server
	port := util.PortBase + 1
	serverIP = net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
	log.Info("Farmer Internal RPC Server using server address", serverIP)
	intCmd = gorpc.NewTCPServer(serverIP, intRPC.NewHandlerFunc())
	err = intCmd.Start()
	if err != nil {
		log.Critical("Farmer Internal RPC service failed to start: ", err)
	}
	defer intCmd.Stop()

	// Tell the world we are alive
	log.Trace("Announcing ourselves to the world")
	go announce()

	// Wait until told told to stop
	<-c
	log.Info("Farmer shutting down")

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Color("red")
	s.Reverse()

	log.Console("Trying to stop Wallet...")
	s.Start()
	select {
	case walletCmd <- "stop":
		s.Stop()
	case <-time.After(5 * time.Second):
		s.Stop()
		log.Console("Wallet timed out - probably wasn't running.")
	}
	log.Console("Trying to stop Databases...")
	s.Start()
	select {
	case dbCmd <- "stop":
		s.Stop()
	case <-time.After(5 * time.Second):
		s.Stop()
		log.Console("Database(s) timed out - probably were not running.")
	}
}
