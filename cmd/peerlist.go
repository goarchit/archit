// Peerlist command - Provides a formated list of all peers
// Originally work created on 5/7/2017
//

package cmd

import (
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/log"
	"github.com/goarchit/archit/util"
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
)

type PeerlistCommand struct {
	PortBase int  `short:"B" long:"PortBase" description:"Primary port number of Archit serve   rs status is requested of" default:"1958" env:"ARCHIT_PORT"`
	SortByRep bool `short:"S" long:"SortByRep" description:"Displays sorted by Reputation instead of the default sorting by IP Addr"`
}

var peerlistCmd PeerlistCommand

func init() {
	_,err := config.Parser.AddCommand("peerlist", "Provides a formated list of peer IPs and their reputations", "", &peerlistCmd)
        if err != nil {
                fmt.Println("Internal error parsing PeerList command:",err)
                os.Exit(1)
        }
}

func (ec *PeerlistCommand) Execute(args []string) error {

	util.PortBase = peerlistCmd.PortBase
	config.Conf(false)

	port := util.PortBase +1
	serverIP := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
	pl := util.GetPeerInfo(serverIP)
	if !peerlistCmd.SortByRep {
		log.Console("Current",len(pl),"peers by IP address:")
		spl := SortPlByIP(pl)
		for i := range spl {
			log.Console("Peer",spl[i].IPAddr,"Rep:",spl[i].Reputation)
		}
	} else {
		log.Console("Current",len(pl),"peers by reputation:")
		spl := util.SortPl(pl)
		for i := range spl {
			log.Console("Peer",spl[i].IPAddr,"Rep:",spl[i].Reputation)
		}
	}

	return nil
}

type sortedPeer struct {
	Reputation int64  // sort key
	WalletAddr string
	IPAddr	string
	}

type ByIP []sortedPeer

func (a ByIP) Len() int 		{ return len(a) }
func (a ByIP) Swap(i, j int) 		{ a[i], a[j] = a[j], a[i] }
func (a ByIP) Less(i, j int) bool 	{ return a[i].IPAddr < a[j].IPAddr }

func SortPlByIP(pl util.PeerList) ByIP {
	sPl := make(ByIP,len(pl))
	// build it...
	i := 0
	for k,v := range pl {
		sPl[i].Reputation = v.Reputation
		sPl[i].WalletAddr = k
		sPl[i].IPAddr = v.IPAddr
		i++
	}

	sort.Sort(ByIP(sPl))
	return sPl
} 
