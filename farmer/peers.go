package farmer

import (
//	"github.com/goarchit/archit/log"
//	"github.com/valyala/gorpc"
)

type Peer struct {
        IpAddr string // including port
        MacAddr string
        Repuatation int64
}

type PeerList struct {
        Address map[string]Peer
}

type PeerInfo struct {
	Address string
	Detail Peer
}

func PeerAdd(pi *PeerInfo) error {
	return nil
}

func PeerDelete(pi *PeerInfo) error {
	return nil
}

func PeerListAll() *PeerList {
	pl := new(PeerList)
	return pl
}
