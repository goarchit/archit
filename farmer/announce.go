package farmer

import (
	"github.com/goarchit/archit/log"
	"github.com/goarchit/archit/util"
	"github.com/valyala/gorpc"
	"encoding/json"
	"net"
	"sync"
	"strings"
)

const SeedPortBase string = ":1958"
var FoundSeed bool
var FarmerMutex sync.Mutex

func announce() {
	
	addresses, err := net.LookupHost("dnsseed.goarchit.online")
	if err != nil {
		log.Critical("Failure to lookup dnsseed.goarchit.online")
	}	
	log.Console("DNSseed resolved to", addresses)

	iAm := new(Peer)
	iAm.IPAddr = util.ServerIP
	iAm.MacAddr = "Invalid"
	rifs := util.RoutedInterface("ip", net.FlagUp|net.FlagBroadcast)
        if rifs != nil {
		iAm.MacAddr = rifs.HardwareAddr.String()
        }
	s, _ := json.Marshal(iAm)
	log.Debug("whoAmI:",string(s))
	
	// Find an active seed node
	for i := 0; i < len(addresses); i++ {
		seed := strings.SplitAfter(addresses[i], " ")
		log.Console("Found seed ",seed)
	//	go tellSeed(s, config.Archit.WalletAddr, seed+SeedPortBase)
	}
	log.Console("Farmer node startup complete!")
}

func tellSeed(p Peer, WalletAddr, serverIP string) {
	var pi PeerInfo
 	c := gorpc.NewTCPClient(serverIP)
        c.Start()
        defer c.Stop()

        d := gorpc.NewDispatcher()
        d.AddFunc("PeerAdd", func(pi *PeerInfo) {})
	d.AddFunc("PeerListAll", func() {})
        dc := d.NewFuncClient(c)
	pi.WalletAddr = WalletAddr
	pi.Detail = p
        _, err := dc.Call("PeerAdd", pi)
        if err != nil {
                log.Warning("Accounce to seed",serverIP,"failed:",err)
        } else {
		FoundSeed = true
		pl, err := dc.Call("PeerListAll", nil)
		if err != nil {
			log.Warning("Attempt to get PeerList from seed",serverIP,"failed:",err)
		} else {
			peerListAdd(pl.(PeerList))
		}
	}
}
