//  Announce announces to the network that this node has joined
//  It does so by telling the first seed node it finds active
package farmer

import (
	"encoding/gob"
	"encoding/json"
	"github.com/goarchit/archit/log"
	"github.com/goarchit/archit/util"
	"github.com/valyala/gorpc"
	"bytes"
	"net"
	"sync"
)

var FarmerMutex sync.Mutex

func announce() {
	var newPL PeerList
	log.Console("Entering annouce()")

	util.Dnsseed()

	log.Console("Dnsseed() called")
	iAm := new(PeerInfo)
	iAm.SenderIP = util.PublicIP
	iAm.WalletAddr = util.WalletAddr
	iAm.HopCount = 2	// Once to my seed, once from there
	iAm.IsASeed = util.IAmASeed
	iAm.Detail.IPAddr = util.PublicIP
	iAm.Detail.MacAddr = "Invalid"
	rifs := util.RoutedInterface("ip", net.FlagUp|net.FlagBroadcast)
	if rifs != nil {
		iAm.Detail.MacAddr = rifs.HardwareAddr.String()
	}
	s, _ := json.Marshal(iAm)
	log.Console("iAm complete")
	log.Debug("iAm:", string(s))

	// Active seed node already found in util.DNSsed() and stored in util.MyDNSServerIP

	c := gorpc.NewTCPClient(util.MyDNSServerIP)
	c.Start()
	defer c.Stop()

	d := gorpc.NewDispatcher()
	d.AddFunc("PeerAdd", func() {})
        d.AddFunc("PeerListAll", func() {})

	dc := d.NewFuncClient(c)
	// Add yourself to your seed node, the seed will tell everyone else
	log.Console("Telling",util.MyDNSServerIP,"that we are",iAm.WalletAddr)
	_, err := dc.Call("PeerAdd", iAm)
	if err != nil {
		log.Critical("Accounce to seed", util.MyDNSServerIP, "failed:", err)
	}
	// Now go ask for all known nodes
	log.Console("Calling PeerListAll at seed",util.MyDNSServerIP)
	plstr, err := dc.Call("PeerListAll",nil)
	if err != nil {
		log.Critical("Attempt to get PeerList from seed", util.MyDNSServerIP, "failed:", err)
	}
	str, ok := plstr.(string)
	if !ok {
		log.Critical("Tellseed: dc.call(PeerListAll) did not return a string")
	}
	buf := bytes.NewBufferString(str)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(&newPL)
	if err != nil {
		log.Critical("Tellseed:  Error decoding PeerMap:",err)
	}
	peerListAdd(newPL)
	log.Trace("Farmer node startup complete!")
}

func tellNode(pi *PeerInfo) {

	serverIP := pi.Detail.IPAddr
	c := gorpc.NewTCPClient(serverIP)
	c.Start()
	defer c.Stop()

	d := gorpc.NewDispatcher()
	// Override SerderIP
	pi.SenderIP = util.PublicIP
	d.AddFunc("PeerAdd", func() {})
	dc := d.NewFuncClient(c)
	log.Console("Calling PeerAdd at",serverIP)
	_, err := dc.Call("PeerAdd", pi)
	if err != nil {
		log.Warning("Announce to node", serverIP, "failed:", err)
	}
}
