// ArchIt Wallet handling  routine
// Originally work created on 1/23/2017
//

package farmer

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcrpcclient"
	"github.com/btcsuite/btcutil"
	"github.com/goarchit/archit/log"
	"github.com/goarchit/archit/util"
	"net"
	"strconv"
)

var imacredit = chaincfg.Params{
	Name:             "IMACNet",
	Net:              0x4d494341,
	PubKeyHashAddrID: 20,
	ScriptHashAddrID: 5,
	HDPrivateKeyID:   [4]byte{0x01, 0x02, 0x03, 0x04},
	HDPublicKeyID:    [4]byte{0x05, 0x06, 0x07, 0x08},
	DefaultPort:      "64097",
	DNSSeeds: []chaincfg.DNSSeed{
		{"dnsseed.imacredit.org", false},
	},
	CoinbaseMaturity:   100,
	TargetTimespan:     256 * 128,
	TargetTimePerBlock: 256,
}

var client *btcrpcclient.Client

func init() {
	err := chaincfg.Register(&imacredit)
	if err != nil {
		log.Critical("Server Chaincfg Register error", err)
	}
}

func Wallet(c chan string, check bool) {
	var err error

	// Start by opening the wallet

	// Connect to local bitcoin core RPC server using HTTP POST mode.
	host := net.JoinHostPort(util.WalletIP, strconv.Itoa(util.WalletPort))

	log.Debug("Host: ", host)
	log.Debug("User: ", util.WalletUser)
	//log.Debug("Password: '",util.WalletPassword,"'")
	connCfg := &btcrpcclient.ConnConfig{
		Host:         host,
		User:         util.WalletUser,
		Pass:         util.WalletPassword,
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}
	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err = btcrpcclient.New(connCfg, nil)
	if err != nil {
		log.Debug("Critical Error connecting to btcclient")
		log.Critical(err)
	}
	defer client.Shutdown()
	if check {
		log.Debug("Specified WalletAddr: ", util.WalletAddr)
		addr, err := btcutil.DecodeAddress(util.WalletAddr, &imacredit)
		if err != nil {
			log.Critical("Wallet error: Invalid WalletAddr specified", err)
		}
		_, err = client.ValidateAddress(addr)
		if err != nil {
			log.Critical("Wallet address validation error:", err, "\n Is the wallet running?")
		} else {
			log.Debug("Wallet address passes validation test")
		}
	}

	// Now go process commands
	cmd := <-c
	for cmd != "stop" {
		switch cmd {
		case "status":
			c <- walletStatus()
		case "generate":
			c <- walletGenerate()
		default:
			log.Critical("Wallet passed unknown command", cmd)
		}
		cmd = <-c
	}
	c <- "Wallet shutdown complete"
}

func walletStatus() string {
	blockCount, err := client.GetBlockCount()
	if err != nil {
		log.Debug("Wallet Error getting blockCount")
		log.Critical(err)
	}
	response := fmt.Sprint("Current Wallet Block count: ", blockCount, "\n")
	connections, err := client.GetConnectionCount()
	if err != nil {
		log.Debug("Wallet Error getting connection count")
		log.Critical(err)
	}
	response += fmt.Sprint("Current Wallet Connections: ", connections, "\n")
	// list accounts
	accounts, err := client.ListAccounts()
	if err != nil {
		log.Critical("Wallet error: ListAccounts", err)
	}
	// iterate over accounts (map[string]btcutil.Amount) and write to stdout
	for label, amount := range accounts {
		response += fmt.Sprint("Account ", label, " balance: ",amount, "\n")
	}
	return response
}

func walletGenerate() string {
	err := client.KeyPoolRefill()
	if err != nil {
		log.Critical("KeyPoolRefill error",err)
	}
	response, err := client.GetNewAddress(util.Account)
	if err != nil {
		log.Debug("Wallet Error getting new address for account",util.Account)
		log.Critical(err)
	}
	log.Debug("WalletGenerate response for account",util.Account,":",response)
	return response.String()
}
