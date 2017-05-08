package farmer

import (
	"bytes"
	"encoding/gob"
	"github.com/goarchit/archit/log"
	"github.com/goarchit/archit/util"
	"net"
	"sync"
)

type PeerInfo struct {
	SenderIP   string // including port
	WalletAddr string
	HopCount   int
	IsASeed    bool
	Detail     util.Peer
}

var PeerMap struct {
	mutex sync.Mutex
	PL    util.PeerList
}
var PeerIP map[string]int  // Indexed by HostID
var PeerMac map[string]int // Indexed by Mac Address

var mapMutex sync.Mutex

const MaxIPsOrMacs int = 50

func init() {
	PeerMap.PL = make(map[string]util.Peer)
	PeerIP = make(map[string]int)
	PeerMac = make(map[string]int)
}

func PeerAdd(pi *PeerInfo) string {

	// Decrement the hop count to prevent if not zero to prevent PeerAdd storms
	if pi.HopCount == 0 {
		return util.OutOfHops
	}
	pi.HopCount--

	// Do quick sanity checks

	if len(pi.WalletAddr) != 34 {
		log.Critical("Invalid WalletAddr:", pi.WalletAddr)
	}
	if pi.WalletAddr[0] != '9' {
		log.Critical("Not an IMAC WalletAddr:", pi.WalletAddr)
	}
	host, _, err := net.SplitHostPort(pi.Detail.IPAddr)
	if err != nil {
		log.Warning("PeerAdd: Received invalid IP address from", pi)
		return err.Error()
	}
	mapMutex.Lock()
	PeerIP[host] += 1
	mapMutex.Unlock()
	if PeerIP[host] > MaxIPsOrMacs {
		return "Too many Farmers behind Public IP " + host
	}
	mapMutex.Lock()
	PeerMac[pi.Detail.MacAddr] += 1
	mapMutex.Unlock()
	if PeerMac[pi.Detail.MacAddr] > MaxIPsOrMacs {
		return "Too many Farmers using MAC Address " + pi.Detail.MacAddr
	}
	if pi.WalletAddr == util.WalletAddr {
		return "" //  Just save outself a lot of processing when other seeds
	} //  tell us about outself
	if pi.IsASeed {
		return "" //  Seeds are not peers and cannot handle data request
	}

	// Onward to the real processing

	val, found := PeerMap.PL[pi.WalletAddr]
	if found {
		log.Info(pi.WalletAddr, "entering network per", pi.SenderIP)
		pm := PeerMap.PL[pi.WalletAddr]
		change := false
		if val.IPAddr != pi.Detail.IPAddr {
			log.Warning("Peer", pi.WalletAddr, "IP has changed!  Old=", val.IPAddr,
				"New=", pi.Detail.IPAddr)
			log.Warning("Accepting change for now, need code to authenticate")
			pm.IPAddr = pi.Detail.IPAddr
			change = true
		}
		if val.MacAddr != pi.Detail.MacAddr {
			log.Warning("Peer", pi.WalletAddr, "MAC has changed!  Old=", val.MacAddr,
				"New=", pi.Detail.MacAddr)
			log.Warning("Accepting change for now, need code to authenticate")
			pm.MacAddr = pi.Detail.MacAddr
			change = true
		}
		if change {
			PeerMap.mutex.Lock()
			PeerMap.PL[pi.WalletAddr] = pm
			PeerMap.mutex.Unlock()
		}
	} else {
		log.Info("Received new peer",pi.WalletAddr,"information from",pi.SenderIP)
		// Only allow the public key to be stored the first time
		peer := PeerMap.PL[pi.WalletAddr]
		peer.PublicKey = pi.Detail.PublicKey
		peer.IPAddr = pi.Detail.IPAddr
		peer.MacAddr = pi.Detail.MacAddr
		peer.Reputation = 0	// Minimum reptuatio
		PeerMap.mutex.Lock()
		PeerMap.PL[pi.WalletAddr] = peer
		PeerMap.mutex.Unlock()
		log.Trace("WalletAddr",pi.WalletAddr,"set to",peer.IPAddr)
	}

	// If your a seed, do your duty and tell everyone
	if pi.HopCount > 0 && util.IAmASeed {
		// Tell the other seed nodes first:$
		for _, v := range util.DNSSeeds {
			ip := v + util.SeedPortBase
			if ip != util.PublicIP {
				go tellNode(pi,ip)
			}
		}

		for _, v := range PeerMap.PL {
			if err != nil {
				log.Critical("Bad news, invalid entry in PeerMap.PL:", err)
			}
			// Don't tell the person that told me
			if v.IPAddr != pi.SenderIP {
				log.Debug("We", util.PublicIP, "are about to tellNode",
					v.IPAddr, "that we learned about", pi.Detail.IPAddr,
					"from", pi.SenderIP)
				go tellNode(pi,v.IPAddr)
			}
		}
	}
	return ""
}

func PeerDelete(key string) error {
	if key == "" {
		log.Critical("PeerDelete called with empty key")
	}
	log.Info("Deleting farmer",key,"at",PeerMap.PL[key].IPAddr,"from PeerMap")
	PeerMap.mutex.Lock()
	delete(PeerMap.PL, key)
	PeerMap.mutex.Unlock()
	return nil
}

func PeerListExt() string {
	if !util.IAmASeed {
		return "Non-Seed nodes will not return a valid Peers List"
	}
	return PeerListAll()
}

func PeerListAll() string {
	var encBuf bytes.Buffer
	enc := gob.NewEncoder(&encBuf)
	err := enc.Encode(&PeerMap.PL)
	if err != nil {
		log.Critical("PeerListAll Encode err:", err)
	}
	return encBuf.String()
}

func peerListAdd(pl util.PeerList) {
	for k, v := range pl {
		log.Trace("Bulk request to add", k, v)
		_, found := PeerMap.PL[k]
		if found {
			log.Debug(k, "already in map, ignored from peer")
		} else if k != util.WalletAddr {
			log.Info(k, "added from seed!")
			// Don't allow sender to set initial Reputation
			v.Reputation = 0
			PeerMap.mutex.Lock()
			PeerMap.PL[k] = v
			PeerMap.mutex.Unlock()
		}
	}
}
