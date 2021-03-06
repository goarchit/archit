// Create command:  Interactively create the configuration file
// Originally work created on 2/4/2017
//

package cmd

import (
	"errors"
	"fmt"
	"github.com/cmiceli/password-generator-go"
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/farmer"
	"github.com/goarchit/archit/log"
	"github.com/goarchit/archit/util"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type CreateCommand struct {
	Account string `short:"A" long:"Account" description:"Account name to use when talking with IMACredit wallet to generate a WalletAddr" default:"Archit"`
	//        Conf    string `short:"c" long:"Conf" description:"Name of the conf file" required:"Y"`
	//       DBDir   string `short:"b" long:"DBDir" description:"Path for persistance databases" required:"Y"`
	Raptor         int    `short:"R" long:"Raptor" description:"Raptor factor - how many extra fountain blocks to generate (4-12)" default:"8" env:"ARCHIT_RAPTOR" choice:"4" choice:"5" choice:"6" choice:"7" choice:"8" choice:"9" choice:"10" choice:"11" choice:"12"`
	PortBase       int    `short:"B" long:"PortBase" description:"Primary port number Archit servers will listen to. Port# +1 will be used interally for server communication" default:"1958"`
	WalletIP       string `long:"WalletIP" short:"I" description:"IP name or address of IMACredit Wallet.  Recommend this be set in your archit configuration file" default:"localhost" env:"ARCHIT_WALLETIP"`
	WalletPort     int    `long:"WalletPort" short:"P" description:"IMACredit Wallets's RPCPort setting." default:"64096" env:"ARCHIT_WALLETPORT"`
	WalletUser     string `long:"WalletUser" short:"U" description:"IMACredit Wallet's RPCuser setting.  Recommend this be set in your archit configuration file" required:"Y"`
	WalletPassword string `long:"WalletPassword" short:"W" description:"IMACredit Wallet's RPCPassword settting.  HIGHLY recommend this be set in your archit configuration file" required:"Y"`
	Pretend        bool   `short:"p" long:"Pretend" description:"Display results instead of actually creating the configuration file"`
	Silent	bool	`short:"s" long:"Silent" description:"Do not display contents of configuration file while creating"`
        MinFreeSpace   int    `long:"MinFreeSpace" short:"M" default:"256" description:"Minimum     free space in GB to maintain in data directory" env:"ARCHIT_MINFREESPACE"`
        DataDir        string `long:"DataDir" short:"D" default:"~/ArchitData" description:"Dire    ctory farmed data should be stored in" env:"ARCHIT_DATADIR"`

}

var createCmd CreateCommand

func init() {
	_, err := config.Parser.AddCommand("create", "Create a new conf file based on passed settings", "", &createCmd)
	if err != nil {
		fmt.Println("Internal error parsing Create command:", err)
		os.Exit(1)
	}
}

func (ec *CreateCommand) Execute(args []string) error {
	var walletCmd = make(chan string)
	var f *os.File

	// Unlike most command with parameters, Create never call Config() to parse the 
	// configuration file.  This makes sense, since its primary usage is to create
	// just that configuartion file.  As such, there is no need to assign most of the
	// util.xyz variables.  

	util.Conf = util.FullPath(config.Archit.Conf)
	util.KeyPass = pwordgen.NewPassword(80)
	// Set for Wallet usage
	util.Account = createCmd.Account  
	util.WalletIP = createCmd.WalletIP
	util.WalletPort = createCmd.WalletPort
	util.WalletUser = createCmd.WalletUser
	util.WalletPassword = createCmd.WalletPassword

	if createCmd.Pretend {
		log.Console("Pretending!  Would create", util.Conf)
	} else {
		log.Console("Generating configuration file", util.Conf)
		basedir := filepath.Dir(util.Conf)
		err := os.MkdirAll(basedir, 0700)
		if err != nil {
			log.Critical("Error creating", basedir, ":", err)
		}

		f, err = os.Create(util.Conf)
		if err != nil {
			log.Critical(err)
		}
		defer f.Close()
	}
	write(f, "Account = "+util.Account)
	write(f, "DBDir = "+config.Archit.DBDir)
	write(f, "DataDir = "+createCmd.DataDir)
	write(f, "MinFreeSpace = "+strconv.Itoa(createCmd.MinFreeSpace))
	write(f, "Raptor = "+strconv.Itoa(createCmd.Raptor))
	write(f, "LogFile = "+config.Archit.LogFile)
	write(f, "LogLevel = "+strconv.Itoa(config.Archit.LogLevel))
	write(f, "ResetLog = "+strconv.FormatBool(config.Archit.ResetLog))
	write(f, "Verbose = "+strconv.Itoa(config.Archit.Verbose))
	write(f, "PortBase = "+strconv.Itoa(createCmd.PortBase))
	write(f, "KeyPass = "+util.KeyPass)
	write(f, "WalletIP = "+util.WalletIP)
	write(f, "WalletPort = "+strconv.Itoa(util.WalletPort))
	write(f, "WalletUser = "+util.WalletUser)
	write(f, "WalletPassword = "+util.WalletPassword)

	if createCmd.Pretend {
		log.Console("WalletAddr = <autogenerated address>")
	} else {
		// OK, we need to cheat here and launch a wallet:
		// Since we creating the configuration file, its a pretty good bet we 
		// dont have a farming node up and running...
		// Fortunately, the farming nodes codebase will work just fine for this

		go farmer.Wallet(walletCmd, false)

		// Tell Wallet server to generate an address
		select {
		case walletCmd <- "generate":
		case <-time.After(10 * time.Second):
			log.Console("IMACredit Wallet timed out - probably wasn't running.")
			return errors.New("Create complete except for WalletAddr =")
		}
		util.WalletAddr = <-walletCmd
		write(f, "WalletAddr = "+util.WalletAddr)

		// Try a clean shutdown
		select {
		case walletCmd <- "stop":
		case <-time.After(5 * time.Second):
			log.Console("Odd, IMACredit Wallet timed out - probably wasn't running.")
		}

	}
	return nil
}

func write(f *os.File, str string) {
	if createCmd.Pretend {
		log.Console(str)
	} else {
		f.Write([]byte(str+"\n"))
		if !createCmd.Silent {
			log.Console(str)
		}
	}
}
