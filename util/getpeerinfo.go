//  GetPeerInfo queries the local farming node and returns a PeerList

package util

import (
	"bytes"
	"encoding/gob"
	"github.com/goarchit/archit/log"
	"github.com/valyala/gorpc"
)

func GetPeerInfo(serverIP string) PeerList {
	var newPL PeerList

        c := gorpc.NewTCPClient(serverIP)
	c.Start()
	defer c.Stop()

	d := gorpc.NewDispatcher()
	d.AddFunc("PeerListAll", func() {})

	dc := d.NewFuncClient(c)

	plstr, err := dc.Call("PeerListAll", nil)
	if err != nil {
		log.Critical("Attempt to get PeerList from local farm node",serverIP,"failed:", err)
	}
	str, ok := plstr.(string)
	if !ok {
		log.Critical("GetPeerInfo: dc.call(PeerListAll) did not return a string")
	}
	buf := bytes.NewBufferString(str)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(&newPL)
	if err != nil {
		log.Critical("GetPeerInfo:  Error decoding PeerMap:", err)
	}
	return newPL
}
