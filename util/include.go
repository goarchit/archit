// Include... common definitions and external variables
// Originally work created on 3/19/2017
//

package util

import (
	"errors"
	"sync"
)

const MaxRaptor int = 12
const ShardLen int = 1024
const SeedPortBase string = ":1958"

var WG sync.WaitGroup
var Mutex sync.Mutex
var SliceName string
var Raptor int
var DerivedKey []byte
var FarmerStop chan bool
var PublicIP string
var DNSSeeds []string
var MyDNSServerIP string
var IAmASeed bool
var OutOfHops = errors.New("Out of Hops")

func init() {
	FarmerStop = make(chan bool)
}
