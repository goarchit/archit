// Include... common definitions and external variables
// Originally work created on 3/19/2017
//

package util

import (
	"sync"
)

const MaxRaptor int = 12
const ShardLen int = 1024

var WG sync.WaitGroup
var SliceName string
var Raptor int
var DerivedKey []byte
var ServerIP string
var FarmerStop chan bool
var PublicIP string
var DNSSeeds []string
var IAmASeed bool

func init() {
	FarmerStop = make(chan bool)
}
