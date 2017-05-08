// Sortpl function: sorts a PeersList
// Originally work created on 5/7/2017
//

package util

import (
	"sort"
)

type ByRep []SortedPeer

func (a ByRep) Len() int 		{ return len(a) }
func (a ByRep) Swap(i, j int) 		{ a[i], a[j] = a[j], a[i] }
func (a ByRep) Less(i, j int) bool 	{ return a[i].Reputation < a[j].Reputation }

func SortPl(pl PeerList) ByRep {
	sPl := make(ByRep,len(pl))	
	// build it...
	i := 0
	for k,v := range pl {
		sPl[i].Reputation = v.Reputation
		sPl[i].WalletAddr = k
		sPl[i].IPAddr = v.IPAddr	
		i++
	}

	sort.Sort(ByRep(sPl))
	return sPl
} 
