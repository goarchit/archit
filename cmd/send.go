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
	"os"
	"time"
)

type SendCommand struct {
 	Chaos 	bool   `short:"C" long:"Chaos" description:"Add a check for a random bit flip (chaos bit) after encoding and decoding - this doubles the computational load for sending data    and normally is seldom needed" env:"ARCHIT_CHAOS"`
	Raptor  int    `short:"R" long:"Raptor" description:"Raptor factor - how many extra    fountain blocks to generate (4-12)" default:"8" env:"ARCHIT_RAPTOR" choice:"4" choice:"5" choice:"6" choice:"7" choice:"8" choice:"9" choice:"10" choice:"11" choice:"12"`
	KeyPass string `short:"k" long:"KeyPass" description:"Your Key Passphase.  Recommend this be set in your archit configuration file" default:"insecure" env:"ARCHIT_KEYPASS"`
        KeyPIN  int    `short:"n" long:"KeyPIN" description:"Your personal identificatio number, used to encrypt KeyPass" default:"0" env:"ARCHIT_KEYPIN"`

}

var sendCmd SendCommand

func init() {
	config.Parser.AddCommand("send", "Send a file [Fee!]", 
		"Sends a file to the Archit network.  This cost IMAC to perform.  Use -C to add extra checks (rarely needed)", &sendCmd)
}

func (ec *SendCommand) Execute(args []string) error {
	var block [32 * util.ShardLen]byte
	var tblock [(32 + util.MaxRaptor) * util.ShardLen]byte
	var iv [32 + util.MaxRaptor][aes.BlockSize]byte

	log.Console("Starting Send Command")
	config.Conf(true) // Get the PIN and derive the key
	db.Open()
	defer db.Close()
	log.Trace("Database open")

	filename := "LICENSE.md"
	db.FileInfo.Filename = []byte(filename)
	db.FileInfo.UploadTime = time.Now()
	f, err := os.Open(filename)
	if err != nil {
		log.Critical(err)
	} else {
		log.Trace("Reading", filename)
	}
	n1, err := f.Read(block[:])
	if err != nil {
		log.Critical(err)
	}
	log.Trace("Read:", n1, "bytes read")
	// Initial tblock
	copy(tblock[:32*util.ShardLen-1], string(block[:]))
	copy(tblock[32*util.ShardLen:32*util.ShardLen+32*config.Archit.Raptor], string(block[:]))

	log.Trace("Starting DBRecord.Slice key determination")
	hash, err := util.HashString(string(block[:]))
	if err != nil {
		log.Critical(err)
	}
	key := fmt.Sprintf("%x", hash)
	log.Trace("Filename:", key)
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
	s.Raptor = config.Archit.Raptor
	db.FileInfo.Slices[util.SliceName] = s
	db.FileInfo.Mutex.Unlock()
	// DB File Information updates must be syncronouse since db.FileInfo is static
	db.FileUpdate()

	log.Console("Sending complete (not really, this is still test code working up to that}")

	return nil
}
