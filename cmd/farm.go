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
	"os"
	"os/signal"
	"strconv"
)

type FarmCommand struct {
        PortBase       int    `short:"B" long:"PortBase" description:"Primary port number Archit servers will listen to. Port# +1 will be used interally for server communication" default:"1958" env:"ARCHIT_PORT"`
        WalletAddr     string `long:"WalletAddr" short:"A" description:"IMACredit Adddress for identity and transactions - !!!Do not use the default!!!" default:"9KsqKMgLjzBWCidAAdy3pNn6jwbd9BT4Te" env:"ARCHIT_WALLETADDR"`
        WalletIP       string `long:"WalletIP" short:"I" description:"IP name or address of IMACredit Wallet.  Recommend this be set in your archit configuration file" default:"localhost" env:"ARCHIT_WALLETIP"`
        WalletPort     int    `long:"WalletPort" short:"P" description:"IMACredit Wallets's RPCPort setting." default:"64096" env:"ARCHIT_WALLETPORT"`
        WalletUser     string `long:"WalletUser" short:"U" description:"IMACredit Wallet's RPCuser setting.  Recommend this be set in your archit configuration file" default:"ReplaceThis"`
        WalletPassword string `long:"WalletPassword" short:"W" description:"IMACredit Wallet's RPCPassword settting.  HIGHLY recommend this be set in your archit configuration file" default:"AbsolutelyChangeThis"`
}

var farmCmd FarmCommand

func init() {
	config.Parser.AddCommand("farm", "Start farming and earning IMAC[earn]", "Starts the ArchIt farming & wallet servers[Earn]", &farmCmd)
}

func (ec *FarmCommand) Execute(args []string) error {
	config.Conf(false)
	pid := os.Getpid()
	err := ioutil.WriteFile(config.Archit.DBDir+"/archit.pid", []byte(strconv.Itoa(pid)), 0644)
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
