// Start command:  Start internal services
// Originally work created on 1/8/2017
//

package cmd

import (
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/farmer"
	"github.com/goarchit/archit/log"
	"github.com/goarchit/archit/util"
	"io/ioutil"
	"fmt"
	"os"
	"os/signal"
	"strconv"
)

type FarmCommand struct {
	SeedMode	bool  `short:"S" long:"SeedMode" description:"Set Seed mode and bypass some checks & Wallet activity" enn:"ARCHIT_SEED"`
        PortBase       int    `short:"B" long:"PortBase" description:"Primary port number Archit servers will listen to. Port# +1 will be used interally for server communication" default:"1958" env:"ARCHIT_PORT"`
        WalletAddr     string `long:"WalletAddr" short:"A" description:"IMACredit Adddress for identity and transactions - !!!Do not use the default!!!" default:"9Ks***INVALID***Ady3pNn6jwbd9BT4Te" env:"ARCHIT_WALLETADDR"`
        WalletIP       string `long:"WalletIP" short:"I" description:"IP name or address of IMACredit Wallet.  Recommend this be set in your archit configuration file" default:"localhost" env:"ARCHIT_WALLETIP"`
        WalletPort     int    `long:"WalletPort" short:"P" description:"IMACredit Wallets's RPCPort setting." default:"64096" env:"ARCHIT_WALLETPORT"`
        WalletUser     string `long:"WalletUser" short:"U" description:"IMACredit Wallet's RPCuser setting.  Recommend this be set in your archit configuration file" default:"ReplaceThis"`
        WalletPassword string `long:"WalletPassword" short:"W" description:"IMACredit Wallet's RPCPassword settting.  HIGHLY recommend this be set in your archit configuration file" default:"AbsolutelyChangeThis"`
}

var farmCmd FarmCommand

func init() {
	_,err := config.Parser.AddCommand("farm", "Start farming and earning IMAC[earn]", "Starts the ArchIt farming & wallet servers[Earn]", &farmCmd)
        if err != nil {
                fmt.Println("Internal error parsing Farm command:",err)
                os.Exit(1)
        }
}

func (ec *FarmCommand) Execute(args []string) error {

	util.SeedMode = farmCmd.SeedMode
	util.PortBase = farmCmd.PortBase
	util.WalletAddr = farmCmd.WalletAddr
	util.WalletIP = farmCmd.WalletIP
	util.WalletPort = farmCmd.WalletPort
	util.WalletUser = farmCmd.WalletUser
	util.WalletPassword = farmCmd.WalletPassword

	if farmCmd.SeedMode {
		log.Console("Seed mode requested.  This means:")
		log.Console("  1) No Wallet information or Keypass required.")
		log.Console("  2) No data storage from renters")
		log.Console("  3) q.e.d.  No earning potential for this node...")
		log.Console("  4) Your PortBase needs to be 1958")
		log.Console("  5) If this is not a registered DNSSeed, this is a waste of effort")
		log.Console("FYI:  Registered seeds are defined in dnsseed.goarchit.online\n")
		util.SeedMode = true
	}
	config.Conf(false) 
	pid := os.Getpid()
	err := ioutil.WriteFile(util.DBDir+"/archit.pid", []byte(strconv.Itoa(pid)), 0644)
	if err != nil {
		log.Critical(err)
	}
	log.Console("Starting Farmer node... Ctrl-C to stop or kill pid", pid)
	go farmer.Run(util.FarmerStop) //  Loops forever waiting on incomming web request
	//  Wait for a Ctrl-C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Console("\nTrying a clean shutdown")
	util.FarmerStop <- true
	return nil
}
