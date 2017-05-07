// Status command:  Display Network and Wallet status
// Originally work created on 1/8/2017
//

package cmd

import (
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/db"
	"github.com/goarchit/archit/log"
	"github.com/goarchit/archit/util"
	"github.com/valyala/gorpc"
	"fmt"
	"os"
	"net"
	"strconv"
)

type StatusCommand struct {
        PortBase int  `short:"B" long:"PortBase" description:"Primary port number Archit s   ervers will listen to. Port# +1 will be used interally for server communication" default:"   1958" env:"ARCHIT_PORT"`
        WalletAddr     string `long:"WalletAddr" short:"A" description:"IMACredit Adddress for identity and transactions - !!!Do not use the default!!!" default:"9Ks***INVALID***Ady3pNn6jwbd9BT4Te" env:"ARCHIT_WALLETADDR"`
        WalletIP       string `long:"WalletIP" short:"I" description:"IP name or address of IMACredit Wallet.  Recommend this be set in your archit configuration file" default:"localhost" env:"ARCHIT_WALLETIP"`
        WalletPort     int    `long:"WalletPort" short:"P" description:"IMACredit Wallets's RPCPort setting." default:"64096" env:"ARCHIT_WALLETPORT"`
        WalletUser     string `long:"WalletUser" short:"U" description:"IMACredit Wallet's RPCuser setting.  Recommend this be set in your archit configuration file" default:"ReplaceThis"`
        WalletPassword string `long:"WalletPassword" short:"W" description:"IMACredit Wallet's RPCPassword settting.  HIGHLY recommend this be set in your archit configuration file" default:"AbsolutelyChangeThis"`

}

var statusCmd StatusCommand

func init() {
	_,err := config.Parser.AddCommand("status", "Shows the status of the ArchIt farming & wallet servers", "", &statusCmd)
        if err != nil {
                fmt.Println("Internal error parsing Status command:",err)
                os.Exit(1)
        }
}

func (ec *StatusCommand) Execute(args []string) error {
	util.PortBase = statusCmd.PortBase
	util.WalletAddr = statusCmd.WalletAddr
	util.WalletIP = statusCmd.WalletIP
	util.WalletPort = statusCmd.WalletPort
	util.WalletUser = statusCmd.WalletUser
	util.WalletPassword = statusCmd.WalletPassword
	config.Conf(false)

	port := util.PortBase + 1
	serverIP := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
	c := gorpc.NewTCPClient(serverIP)
	c.Start()
	defer c.Stop()

	d := gorpc.NewDispatcher()
	d.AddFunc("Status", func() {})
	dc := d.NewFuncClient(c)
	response, err := dc.Call("Status", nil)
	if err != nil {
		log.Error("Status failed: ", err)
	}
	log.Console(response)

	// Display local DB Stats
        db.Open()
        defer db.Close()
        s := db.Status()
        log.Console(s)

	return nil
}
