package farmer

import (
	"encoding/gob"
	"encoding/json"
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/log"
	"github.com/goarchit/archit/util"
	"github.com/valyala/gorpc"
	"bytes"
	"net"
	"sync"
)

const SeedPortBase string = ":1958"

var FoundSeed bool
var FarmerMutex sync.Mutex

func announce() {

	addresses, err := net.LookupHost("dnsseed.goarchit.online")
	if err != nil {
		log.Critical("Failure to lookup dnsseed.goarchit.online")
	}
	log.Debug("DNSseed resolved to", addresses)

	iAm := new(PeerInfo)
	iAm.WalletAddr = config.Archit.WalletAddr
	iAm.Detail.IPAddr = util.ServerIP
	iAm.Detail.MacAddr = "Invalid"
	rifs := util.RoutedInterface("ip", net.FlagUp|net.FlagBroadcast)
	if rifs != nil {
		iAm.Detail.MacAddr = rifs.HardwareAddr.String()
	}
	s, _ := json.Marshal(iAm)
	log.Debug("whoAmI:", string(s))

	// Find an active seed node
	for i := 0; i < len(addresses); i++ {
		seed := addresses[i]
		log.Trace("Found seed ", seed)
		go tellSeed(iAm,seed+SeedPortBase)
	}
	log.Trace("Farmer node startup complete!")
}

func tellSeed(pi *PeerInfo,serverIP string) {

	var newPL PeerList

	c := gorpc.NewTCPClient(serverIP)
	c.Start()
	defer c.Stop()

	d := gorpc.NewDispatcher()
	d.AddFunc("PeerAdd", func(pi *PeerInfo) {})
	d.AddFunc("PeerListAll", func() {})
	dc := d.NewFuncClient(c)
	_, err := dc.Call("PeerAdd", pi)
	if err != nil {
		log.Warning("Accounce to seed", serverIP, "failed:", err)
	} else {
		FoundSeed = true
		s, err := dc.Call("PeerListAll",nil)
		if err != nil {
			log.Warning("Attempt to get PeerList from seed", serverIP, "failed:", err)
		} else {
			str, ok := s.(string)
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
		}
	}
}
