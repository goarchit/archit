// Send command:  Send a file to the network
// Originally work created on 2/4/2017
//

package cmd

import (
	"crypto/aes"
	"fmt"
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/db"
	"github.com/goarchit/archit/log"
	"github.com/goarchit/archit/util"
	"math"
	"net"
	"os"
	"strconv"
	"time"
)

type SendCommand struct {
        Account string `short:"A" long:"Account" description:"IMACredit wallet account name to use for payment" default:"Archit"`
 	Chaos 	bool   `short:"C" long:"Chaos" description:"Add a check for a random bit flip (chaos bit) after encoding and decoding - this doubles the computational load for sending data    and normally is seldom needed" env:"ARCHIT_CHAOS"`
	Raptor  int    `short:"R" long:"Raptor" description:"Raptor factor - how many extra    fountain blocks to generate (4-12)" default:"8" env:"ARCHIT_RAPTOR" choice:"4" choice:"5" choice:"6" choice:"7" choice:"8" choice:"9" choice:"10" choice:"11" choice:"12"`
	KeyPass string `short:"k" long:"KeyPass" description:"Your Key Passphase.  Recommend this be set in your archit configuration file" default:"insecure" env:"ARCHIT_KEYPASS"`
        KeyPIN  int    `short:"n" long:"KeyPIN" description:"Your personal identificatio number, used to encrypt KeyPass" default:"0" env:"ARCHIT_KEYPIN"`
	FileName string `short:"F" long:"FileName" description:"File to send.  Will prompt if omitted" default:"" env:"ARCHIT_FILENAME"`
}

var sendCmd SendCommand

func init() {
	_,err := config.Parser.AddCommand("send", "Send a file [Fee!]", 
		"Sends a file to the Archit network.  This cost IMAC to perform.  Use -C to add extra checks (rarely needed)", &sendCmd)
        if err != nil {
                fmt.Println("Internal error parsing Send command:",err)
                os.Exit(1)
        }
}

func (ec *SendCommand) Execute(args []string) error {
	var block [32 * util.ShardLen]byte
	var tblock [(32 + util.MaxRaptor) * util.ShardLen]byte
	var iv [32 + util.MaxRaptor][aes.BlockSize]byte

	util.Account = sendCmd.Account
	util.Chaos = sendCmd.Chaos
	util.Raptor = sendCmd.Raptor
	util.KeyPass = sendCmd.KeyPass
	util.KeyPIN = sendCmd.KeyPIN
	log.Console("KeyPIN =",util.KeyPIN)
	fileName := sendCmd.FileName
	config.Conf(true) // Get the PIN and derive the key
	if sendCmd.Chaos {
		log.Warning("Chaos option ignored, not implemented yet")
	}
	// Go grab the filename if not passed as a parameter
	if sendCmd.FileName == "" {
		fmt.Printf("Name of file to be sent: ")
		_, err := fmt.Scanf("%s",&fileName)
		if err != nil {
			log.Critical("Error entering file name:",err)
		} 
	}
	fileName = util.FullPath(fileName)
	// And check if it exists!
	db.FileInfo.Filename = []byte(fileName)
	db.FileInfo.UploadTime = time.Now()
	f, err := os.Open(fileName)
	if err != nil {
		log.Critical(err)
	} else {
		log.Trace("Attempting to send", fileName)
	}
	fileInfo, err := f.Stat()
	if err != nil {
		log.Critical("Unexpected error Stat()ing file since its already open!",err)
	}
	blockCount := math.Ceil(float64(fileInfo.Size())/float64(util.GB))
	log.Info("About to send",blockCount,"file slices")

	// Grab peers we can send too..
        port := util.PortBase + 1
        serverIP := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
	pl := util.GetPeerInfo(serverIP)
	log.Debug("Sending of",fileName,"Potential Peers:",len(pl))
	if len(pl) < (32+util.Raptor)  {
		log.Critical("Only",len(pl),"Farmers known.  Not enough to function")
	}
	if len(pl) < (42+util.Raptor) {
		log.Warning("Only",len(pl),"Farmers known... sending may fail if at least",32+util.Raptor,"are not online")
	} else {
		if len(pl) < (int(blockCount)*(32+util.Raptor)) {
			log.Warning("Only",len(pl),"Farmers known... some farmers will receive multiple shards (which is ok, just not ideal)")
		}
	}

	log.Trace("UnSorted pl:",pl)
	spl := util.SortPl(pl)
	log.Debug("Sorted pl:",spl)

	db.Open()
	defer db.Close()
	log.Trace("Send command database open")


	n1, err := f.Read(block[:])
	if err != nil {
		log.Critical(err)
	}
	log.Info("Read:", n1, "bytes read from",fileName)
	// Initial tblock
	copy(tblock[:32*util.ShardLen-1], string(block[:]))
	copy(tblock[32*util.ShardLen:32*util.ShardLen+32*util.Raptor], string(block[:]))

	log.Trace("Starting DBRecord.Slice key determination")
	hash, err := util.HashString(string(block[:]))
	if err != nil {
		log.Critical(err)
	}
	key := fmt.Sprintf("%x", hash)
	log.Info("Filename for farmers:", key)
	util.SliceName = key
	// Encode
	log.Trace("Encoding starting")
	util.WG.Add(32)
	for i := 0; i < 32; i++ {
		go util.Encode(&tblock, &block, i)
	}
	util.WG.Wait()
	// Encrypt
	log.Trace("Encrypting starting")
	util.WG.Add(32)
	for i := 0; i < 32; i++ {
		off := i * util.ShardLen
		go util.Encrypt(tblock[off:off+util.ShardLen-1], &iv, i)
	}
	util.WG.Wait()
	// Create Authentication string
	log.Trace("Authentication generations starting")
	util.WG.Add(32)
	for i := 0; i < 32; i++ {
		off := i * util.ShardLen
		go db.HMAC(string(tblock[off:off+util.ShardLen-1]), i)
	}
	util.WG.Wait()
	log.Trace("Housekeeping tasks starting")
	db.FileInfo.Mutex.Lock()
	s := db.FileInfo.Slices[util.SliceName]
	s.Raptor = util.Raptor
	db.FileInfo.Slices[util.SliceName] = s
	db.FileInfo.Mutex.Unlock()
	// DB File Information updates must be syncronouse since db.FileInfo is static
	db.FileUpdate()
	db.Sync()	// Flush new info to disk

	log.Console("Sending complete (not really, this is still test code working up to that}")

	return nil
}
