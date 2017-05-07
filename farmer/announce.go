//  Announce announces to the network that this node has joined
//  It does so by telling the first seed node it finds active
package farmer

import (
	"encoding/json"
	"github.com/goarchit/archit/log"
	"github.com/goarchit/archit/util"
	"github.com/valyala/gorpc"
	"net"
	"sync"
)

var FarmerMutex sync.Mutex

func announce() {
	var newPL util.PeerList
	var iAm PeerInfo

	if util.MyDNSServerIP == "" {
		log.Console("We are alone... so lonely... please start up another Seed node!")
		return
	}

	if !util.IAmASeed {
		iAm.SenderIP = util.PublicIP
		iAm.WalletAddr = util.WalletAddr
		iAm.HopCount = 2 // Once to my seed, once from there
		iAm.IsASeed = util.IAmASeed
		iAm.Detail.IPAddr = util.PublicIP
		iAm.Detail.MacAddr = "Invalid"
		rifs := util.RoutedInterface("ip", net.FlagUp|net.FlagBroadcast)
		if rifs != nil {
			iAm.Detail.MacAddr = rifs.HardwareAddr.String()
		}
		s, _ := json.Marshal(iAm)
		log.Debug("iAm:", string(s))
	}

	// Active seed node already found in util.DNSsed() and stored in util.MyDNSServerIP

	c := gorpc.NewTCPClient(util.MyDNSServerIP)
	c.Start()
	defer c.Stop()

	d := gorpc.NewDispatcher()
	d.AddFunc("PeerAdd", func() {})

	dc := d.NewFuncClient(c)
	if !util.IAmASeed {
		// Add yourself to your seed node, the seed will tell everyone else
		log.Debug("Telling", util.MyDNSServerIP, "that we are", iAm.WalletAddr)
		_, err := dc.Call("PeerAdd", iAm)
		if err != nil {
			log.Critical("Announce to seed", util.MyDNSServerIP, "failed:", err)
		}
	}
	// Now go ask for all known nodes
	log.Trace("Calling PeerListAll at seed", util.MyDNSServerIP)
	newPL = util.GetPeerInfo(util.MyDNSServerIP)
	peerListAdd(newPL)
	if util.IAmASeed {
		log.Console("Waiting for farmers to join the network")
	} else {
		log.Console("Farmer node startup complete!")
	}
}

func tellNode(pi *PeerInfo, nodeIP string) {

	c := gorpc.NewTCPClient(nodeIP)
	c.Start()
	defer c.Stop()

	d := gorpc.NewDispatcher()
	// Override SerderIP
	pi.SenderIP = util.PublicIP
	d.AddFunc("PeerAdd", func() {})
	dc := d.NewFuncClient(c)
	log.Debug("Calling PeerAdd at", nodeIP)
	_, err := dc.Call("PeerAdd", pi)
	if err != nil {
		log.Warning("Announce to node", nodeIP, "failed:", err)
	}
}
