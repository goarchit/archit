//  Checkpeers() is called once an hour to Ping all members in the PeerList and:
//  1) Adjust reputation up 1 if they respond
//  2) Adjust reputation down 5 if they do not
//  3) Remove them from PeerList if reputation has gone negatrive

package farmer

import (
	"github.com/goarchit/archit/log"
	"github.com/goarchit/archit/util"
	"github.com/valyala/gorpc"

)

func CheckPeers() {
	for k,v := range PeerMap.PL {
		go checkPeer(k,v)
	}
}

func checkPeer(key string,peer util.Peer) {
	theirWalletAddr := CheckWallet(peer.IPAddr)
	if key != theirWalletAddr {
		PeerDelete(key)
		return
	}
	if theirWalletAddr != "" {
		PeerMap.mutex.Lock()
		peer := PeerMap.PL[key]
		peer.Reputation++
		PeerMap.PL[key] = peer
		PeerMap.mutex.Unlock()
		log.Trace("Reputation of",peer.IPAddr,"improved by 1")
	} else {
		PeerMap.mutex.Lock()
		peer := PeerMap.PL[key]
		peer.Reputation -= 5
		PeerMap.PL[key] = peer
		PeerMap.mutex.Unlock()
		log.Trace("Reputation of",peer.IPAddr,"decreased by 5")
		if PeerMap.PL[key].Reputation < 0 {
			PeerDelete(key)
		}
	}
}

func CheckWallet(serverIP string) string {
	c := gorpc.NewTCPClient(serverIP)
	c.Start()
	defer c.Stop()

        d := gorpc.NewDispatcher()
        d.AddFunc("WhoAreYou", func() {})
        dc := d.NewFuncClient(c)
        res, err := dc.Call("Ping", nil)
        if err != nil {
                log.Error("CheckWallet for",serverIP,"failed: ", err)
		return ""
        }
	str, ok := res.(string)
	if !ok {
		log.Critical("CheckWallet: WhoAreYou didn't return a string")
	}
	return str
}
