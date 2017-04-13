// As an alternative to complaining about a missing configuration file, go generate one!

package config

import (
	"github.com/goarchit/archit/log"
	"github.com/goarchit/archit/util"
	"github.com/cmiceli/password-generator-go"
	"os"
)

func genGoodConfig() {
	// first, out with the old
	log.Info("Error processinga",util.Conf," Building a basic KeyPass conf for you!")
	os.Remove(Archit.Conf)
	f, err := os.Create(Archit.Conf)
	if err != nil {
                log.Critical(err)
        } else {
                log.Trace("Creating Basic", util.Conf)
		
        }
	defer f.Close()
 	f.Write([]byte("KeyPass = "+pwordgen.NewPassword(80)+"\n"))
}
