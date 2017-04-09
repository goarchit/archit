// Status command:  Display Network and Wallet status
// Originally work created on 1/8/2017
//

package cmd

import (
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/db"
	"github.com/goarchit/archit/log"
	"github.com/valyala/gorpc"
	"net"
	"strconv"
)

type StatusCommand struct {
        WalletAddr     string `long:"WalletAddr" short:"A" description:"IMACredit Adddress for identity and transactions - !!!Do not use the default!!!" default:"9KsqKMgLjzBWCidAAdy3pNn6jwbd9BT4Te" env:"ARCHIT_WALLETADDR"`
        WalletIP       string `long:"WalletIP" short:"I" description:"IP name or address of IMACredit Wallet.  Recommend this be set in your archit configuration file" default:"localhost" env:"ARCHIT_WALLETIP"`
        WalletPort     int    `long:"WalletPort" short:"P" description:"IMACredit Wallets's RPCPort setting." default:"64096" env:"ARCHIT_WALLETPORT"`
        WalletUser     string `long:"WalletUser" short:"U" description:"IMACredit Wallet's RPCuser setting.  Recommend this be set in your archit configuration file" default:"ReplaceThis"`
        WalletPassword string `long:"WalletPassword" short:"W" description:"IMACredit Wallet's RPCPassword settting.  HIGHLY recommend this be set in your archit configuration file" default:"AbsolutelyChangeThis"`

}

var statusCmd StatusCommand

func init() {
	config.Parser.AddCommand("status", "Shows the status of the ArchIt farming & wallet servers", "", &statusCmd)
}

func (ec *StatusCommand) Execute(args []string) error {
	config.Conf(false)

	// Display local DB Stats
        db.Open()
        defer db.Close()
        s := db.Status()
        log.Console(s)

	// Insert RPC code to query the farmer

	port := config.Archit.PortBase + 1
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
	return nil
}
