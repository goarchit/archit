// As an alternative to complaining about a missing configuration file, go generate one!

package config

import (
	"github.com/goarchit/archit/log"
	"github.com/cmiceli/password-generator-go"
	"os"
)

func genGoodConfig() {
	// first, out with the old
	log.Console("Error processinga",Archit.Conf," Building a basic file for you!")
	os.Remove(Archit.Conf)
	f, err := os.Create(Archit.Conf)
	if err != nil {
                log.Critical(err)
        } else {
                log.Trace("Creating Basic", Archit.Conf)
		
        }
	defer f.Close()
 	f.Write([]byte("KeyPass = "+pwordgen.NewPassword(80)+"\n"))
	f.Write([]byte("RPCuser = architrpc\n"))
	f.Write([]byte("RPCPassword = "+pwordgen.NewPassword(48)+"\n"))
	log.Console("It is recommended you add WalletAddr, WalletUser & WalletPassword or specify them via the command line or environment variables!")
}
