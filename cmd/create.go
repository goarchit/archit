// Create command:  Interactively create the configuration file
// Originally work created on 2/4/2017
//

package cmd

import (
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/log"
	"github.com/goarchit/archit/util"
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

	newConf := util.Conf+".new"
	os.Remove(newConf)
	f, err := os.Create(newConf)
	if err != nil {
		log.Critical(err)
	} else {
		log.Trace("Creating",newConf)
	}
	defer f.Close()
	if util.KeyPass == "insecure" {
		log.Debug("Keypass oddly defaulted, generating random one")
		util.KeyPass = pwordgen.NewPassword(80)
	}
	f.Write([]byte("KeyPass = "+util.KeyPass+"\n"))
	f.Write([]byte("DBDir = "+util.DBDir+"\n"))
	f.Write([]byte("LogFile = "+util.LogFile+"\n"))
	f.Write([]byte("LogLevel = "+strconv.Itoa(util.LogLevel)+"\n"))
	f.Write([]byte("Verbose = "+strconv.Itoa(util.Verbose)+"\n"))
	if util.ResetLog {
		f.Write([]byte("ResetLog = true\n"))
	} else {
		f.Write([]byte("ResetLog = false\n"))
	}
	f.Write([]byte("Raptor = "+strconv.Itoa(util.Raptor)+"\n"))
	f.Write([]byte("PortBase = "+strconv.Itoa(util.PortBase)+"\n"))
	f.Write([]byte("WalletAddr = "+util.WalletAddr+"\n"))
	f.Write([]byte("WalletIP = "+util.WalletIP+"\n"))
	f.Write([]byte("WalletPort = "+strconv.Itoa(util.WalletPort)+"\n"))
	f.Write([]byte("WalletUser = "+util.WalletUser+"\n"))
	f.Write([]byte("WalletPassword = "+util.WalletPassword+"\n"))

	err = f.Close()
	if err != nil {
		log.Critical(err)
	}

	log.Console("Attempting to remove any previous .old conf file")
	os.Remove(util.Conf+".old")
	log.Console("Renameing old configuration file...")
	os.Rename(util.Conf,util.Conf+".old")
	log.Console("Moving new file into place")
	err = os.Rename(util.Conf+".new",util.Conf)
	if err != nil {
		log.Critical("Something bad happened, your old file is likely stored with the extension .old",err)
	}

	return nil
}
