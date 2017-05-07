// Include... common definitions and external variables
// Originally work created on 3/19/2017
//

package util

import (
	"sync"
)

const GB uint64 = 1024*1024*1024
const MaxRaptor int = 12
const ShardLen int = 1024
const SeedPortBase  = ":1958"
const OutOfHops = "Out of Hops"
const FileDBName = "FileInfo.bolt"
const PeerDBName = "PeerInfo.bolt"

var WG sync.WaitGroup
var Mutex sync.Mutex
var SliceName string
var DerivedKey []byte
var FarmerStop chan bool
var PublicIP string
var DNSSeeds []string
var MyDNSServerIP string
var IAmASeed bool

//  Store all flags here so there are in a common place to ease programming

//  Global command line flags found in ../config/parser.go
//  and likely in the configuration file
var Conf string
var DBDir string
var LogFile string
var LogLevel int
var ResetLog bool
var Verbose int

//  Subcommand flags that may be found in the configuration file
var Account string
var SeedMode bool
var PortBase int
var WalletAddr string
var WalletIP string
var WalletPort int
var WalletUser string
var WalletPassword string
var Chaos bool
var Raptor int
var KeyPass string
var KeyPIN int
var MinFreeSpace uint64
var DataDir string

// Peer structures used in several places

type Peer struct {
        IPAddr     string // including port
        MacAddr    string
        Reputation int64
        PublicKey  string
}

type PeerList map[string]Peer // Indexed by Wallet Address

func init() {
	FarmerStop = make(chan bool)
}
