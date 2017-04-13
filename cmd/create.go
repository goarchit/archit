// Create command:  Interactively create the configuration file
// Originally work created on 2/4/2017
//

package cmd

import (
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/log"
	"github.com/goarchit/archit/server"
	"github.com/goarchit/archit/util"
	"github.com/cmiceli/password-generator-go"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

type CreateCommand struct {
	Account string `short:"A" long:"Account" description:"Account name to use when talking with IMACredit wallet to generate a WalletAddr" default:"Archit"`
        Conf    string `short:"c" long:"Conf" description:"Name of the conf file" required:"Y"`
        DBDir   string `short:"b" long:"DBDir" description:"Path for persistance databases" required:"Y"`
	Raptor  int    `short:"R" long:"Raptor" description:"Raptor factor - how many extra fountain blocks to generate (4-12)" default:"8" env:"ARCHIT_RAPTOR" choice:"4" choice:"5" choice:"6" choice:"7" choice:"8" choice:"9" choice:"10" choice:"11" choice:"12"`
	PortBase       int    `short:"B" long:"PortBase" description:"Primary port number Archit servers will listen to. Port# +1 will be used interally for server communication" default:"1958"`
        WalletIP       string `long:"WalletIP" short:"I" description:"IP name or address of IMACredit Wallet.  Recommend this be set in your archit configuration file" default:"localhost" env:"ARCHIT_WALLETIP"`
        WalletPort     int    `long:"WalletPort" short:"P" description:"IMACredit Wallets's RPCPort setting." default:"64096" env:"ARCHIT_WALLETPORT"`
        WalletUser     string `long:"WalletUser" short:"U" description:"IMACredit Wallet's RPCuser setting.  Recommend this be set in your archit configuration file" required:"Y"`
        WalletPassword string `long:"WalletPassword" short:"W" description:"IMACredit Wallet's RPCPassword settting.  HIGHLY recommend this be set in your archit configuration file" required:"Y"`


}

var createCmd CreateCommand

func init() {
	_,err := config.Parser.AddCommand("create", "Create a new conf file based on passed settings", "", &createCmd)
	if err != nil {
		fmt.Println("Internal error parsing Create command:",err)
		os.Exit(1)
	}
}

func (ec *CreateCommand) Execute(args []string) error {
	var walletCmd = make(chan string)

	util.Conf = createCmd.Conf
	util.DBDir = createCmd.DBDir
	util.Raptor = createCmd.Raptor
	util.PortBase = createCmd.PortBase
	util.WalletIP = createCmd.WalletIP
	util.WalletPort = createCmd.WalletPort
	util.WalletUser = createCmd.WalletUser
	util.WalletPassword = createCmd.WalletPassword

	os.Remove(util.Conf)

	f, err := os.Create(util.Conf)
	if err != nil {
		log.Critical(err)
	} 
	defer f.Close()
	util.KeyPass = pwordgen.NewPassword(80)
	f.Write([]byte("KeyPass = "+util.KeyPass+"\n"))
	f.Write([]byte("DBDir = "+util.DBDir+"\n"))
	f.Write([]byte("Raptor = "+strconv.Itoa(util.Raptor)+"\n"))
	f.Write([]byte("LogFile = "+util.LogFile+"\n"))
	f.Write([]byte("LogLevel = "+strconv.Itoa(util.LogLevel)+"\n"))
	f.Write([]byte("ResetLog = "+strconv.FormatBool(util.ResetLog)+"\n"))
	f.Write([]byte("Verbose = "+strconv.Itoa(util.Verbose)+"\n"))
	f.Write([]byte("PortBase = "+strconv.Itoa(util.PortBase)+"\n"))
	f.Write([]byte("WalletIP = "+util.WalletIP+"\n"))
	f.Write([]byte("WalletPort = "+strconv.Itoa(util.WalletPort)+"\n"))
	f.Write([]byte("WalletUser = "+util.WalletUser+"\n"))
	f.Write([]byte("WalletPassword = "+util.WalletPassword+"\n"))

	go server.Wallet(walletCmd,false)
	
	// Tell Wallet server to generate an address
        select {
        case walletCmd <- "generate":
        case <-time.After(10 * time.Second):
                log.Console("IMACredit Wallet timed out - probably wasn't running.")
		return errors.New("Create complete except for WalletAddr =")
        }
	util.WalletAddr = <- walletCmd

	// Try a clean shutdown
        select {
        case walletCmd <- "stop":
        case <-time.After(5 * time.Second):
                log.Console("Odd, IMACredit Wallet timed out - probably wasn't running.")
        }

	f.Write([]byte("WalletAddr = "+util.WalletAddr+"\n"))
	err = f.Close()
	if err != nil {
		log.Critical(err)
	}

	return nil
}
