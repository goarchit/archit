// Create command:  Interactively create the configuration file
// Originally work created on 2/4/2017
//

package cmd

import (
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/log"
	"github.com/cmiceli/password-generator-go"
	"os"
	"strconv"
)

type CreateCommand struct {
}

var createCmd CreateCommand

func init() {
	config.Parser.AddCommand("create", "Create a new conf file based on current settings", "", &createCmd)
}

func (ec *CreateCommand) Execute(args []string) error {

	config.Conf(false)

	log.Console("Creating configuration file based on current settings")

	newConf := config.Archit.Conf+".new"
	os.Remove(newConf)
	f, err := os.Create(newConf)
	if err != nil {
		log.Critical(err)
	} else {
		log.Trace("Creating",newConf)
	}
	defer f.Close()
	if config.Archit.KeyPass == "insecure" {
		log.Debug("Keypass oddly defaulted, generating random one")
		config.Archit.KeyPass = pwordgen.NewPassword(80)
	}
	f.Write([]byte("KeyPass = "+config.Archit.KeyPass+"\n"))
	f.Write([]byte("DBDir = "+config.Archit.DBDir+"\n"))
	f.Write([]byte("LogFile = "+config.Archit.LogFile+"\n"))
	f.Write([]byte("LogLevel = "+strconv.Itoa(config.Archit.LogLevel)+"\n"))
	f.Write([]byte("Verbose = "+strconv.Itoa(config.Archit.Verbose)+"\n"))
	if config.Archit.ResetLog {
		f.Write([]byte("ResetLog = true\n"))
	} else {
		f.Write([]byte("ResetLog = false\n"))
	}
	f.Write([]byte("Raptor = "+strconv.Itoa(config.Archit.Raptor)+"\n"))
	f.Write([]byte("Portbase = "+strconv.Itoa(config.Archit.PortBase)+"\n"))
	f.Write([]byte("WalletAddr = "+config.Archit.WalletAddr+"\n"))
	f.Write([]byte("WalletIP = "+config.Archit.WalletIP+"\n"))
	f.Write([]byte("WalletPort = "+strconv.Itoa(config.Archit.WalletPort)+"\n"))
	f.Write([]byte("WalletUser = "+config.Archit.WalletUser+"\n"))
	f.Write([]byte("WalletPassword = "+config.Archit.WalletPassword+"\n"))

	err = f.Close()
	if err != nil {
		log.Critical(err)
	}

	log.Console("Attempting to remove any previous .old conf file")
	os.Remove(config.Archit.Conf+".old")
	log.Console("Renameing old configuration file...")
	os.Rename(config.Archit.Conf,config.Archit.Conf+".old")
	log.Console("Moving new file into place")
	err = os.Rename(config.Archit.Conf+".new",config.Archit.Conf)
	if err != nil {
		log.Critical("Something bad happened, your old file is likely stored with the extension .old",err)
	}

	return nil
}
