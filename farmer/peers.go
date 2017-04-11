package farmer

import (
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/log"
	"github.com/goarchit/archit/util"
	"bytes"
	"encoding/gob"
	"net"
	"sync"
)

type Peer struct {
	IPAddr     string // including port
	MacAddr    string
	Reputation int64
	PublicKey  string
}

type PeerInfo struct {
	SenderIP   string // including port
	WalletAddr string
	HopCount   int
	Detail     Peer
}

type PeerList map[string]Peer // Indexed by Wallet Address

var PeerMap struct {
	mutex sync.Mutex
	PL    PeerList
}
var PeerIP map[string]int  // Indexed by HostID
var PeerMac map[string]int // Indexed by Mac Address

const MaxIPsOrMacs int = 50

func init() {
	PeerMap.PL = make(map[string]Peer)
	PeerIP = make(map[string]int)
	PeerMac = make(map[string]int)
}

func PeerAdd(pi *PeerInfo) string {
	PeerMap.mutex.Lock()		// Needed?
	defer PeerMap.mutex.Unlock()

	// Decrement the hop count to prevent if not zero to prevent PeerAdd storms
	if pi.HopCount == 0 {
		return util.OutOfHops
	}
	pi.HopCount--

	// Do quick sanity checks

	if len(pi.WalletAddr) != 34 {
		log.Critical("Invalid WalletAddr:",pi.WalletAddr)
	}
	if pi.WalletAddr[0] != '9' {
		log.Critical("Not an IMAC WalletAddr:",pi.WalletAddr)
	}
	host, _, err := net.SplitHostPort(pi.Detail.IPAddr)
	if err != nil {
		log.Warning("PeerAdd: Received invalid IP address from", pi)
		return err.Error()
	}
	PeerIP[host] += 1
	if PeerIP[host] > MaxIPsOrMacs {
		return "Too many Farmers behind Public IP " + host
	}
	PeerMac[pi.Detail.MacAddr] += 1
	if PeerMac[pi.Detail.MacAddr] > MaxIPsOrMacs {
		return "Too many Farmers using MAC Address " + pi.Detail.MacAddr
	}
	if pi.WalletAddr == config.Archit.WalletAddr {
		return ""  //  Just save outself a lot of processing when other seeds 
	}		   //  tell us about outself
	

	// Onward to the real processing

	val, found := PeerMap.PL[pi.WalletAddr]
	if found {
		log.Console(pi.WalletAddr, "entering network per",pi.SenderIP)
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
			PeerMap.PL[pi.WalletAddr] = pm
		}
	} else {
		log.Console(pi.WalletAddr,"is new to the PeerList")
		// Only allow the public key to be stored the first time
		pm := PeerMap.PL[pi.WalletAddr]
		pm.PublicKey = pi.Detail.PublicKey
		pm.IPAddr = pi.Detail.IPAddr
		pm.MacAddr = pi.Detail.MacAddr
		pm.Reputation = 0
		PeerMap.PL[pi.WalletAddr] = pm
	}

	// If your a seed, do your duty and tell everyone
	if util.IAmASeed {
		for _, v := range PeerMap.PL {
			if err != nil {
				log.Critical("Bad news, invalid entry in PeerMap.PL:",err)
			}
			// Don't tell myself or the person that told me
			if v.IPAddr != util.PublicIP && v.IPAddr != pi.SenderIP{
				log.Debug("We",util.PublicIP,"are about to tellNode",
					v.IPAddr,"that we learned about",pi.Detail.IPAddr,
					"from",pi.SenderIP)
				go tellNode(pi)
			}
		}
	}
	return ""
}

func PeerDelete(p *Peer) error {
	return nil
}

func PeerListAll() string {
	var encBuf bytes.Buffer
	enc := gob.NewEncoder(&encBuf)
	err := enc.Encode(&PeerMap.PL)
	if err != nil {
		log.Critical("PeerListAll Encode err:",err)
	}
	return encBuf.String()
}

func peerListAdd(pl PeerList) {
	PeerMap.mutex.Lock()
	defer PeerMap.mutex.Unlock()
	for k, v := range pl {
		log.Trace("Bulk request to add", k, v)
		_, found := PeerMap.PL[k]
		if found {
			log.Debug(k,"already in map, ignored from peer")
		} else if k != config.Archit.WalletAddr {
			log.Console(k,"added from peer!")
			// Don't allow sender to set initial Reputation
			v.Reputation = 0
			PeerMap.PL[k] = v
		}
	}
}
