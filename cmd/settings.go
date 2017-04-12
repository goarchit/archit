// Settings command:  Display all configurable settings
// Originally work created on 1/8/2017
//

package cmd

import (
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/log"
	"github.com/goarchit/archit/util"
)

type SettingsCommand struct {
	NoHide	bool `short:"s" long:"ShowPasswords" description:"Show passwords" env:"ARCHIT_NOHIDE"`
        //Chaos   bool   `short:"C" long:"Chaos" description:"Add a check for a random bit flip (chaos bit) after encoding and decoding - this doubles the computational load for sending data    and normally is seldom needed" env:"ARCHIT_CHAOS"`
        //Raptor  int    `short:"R" long:"Raptor" description:"Raptor factor - how many extra    fountain blocks to generate (4-12)" default:"8" env:"ARCHIT_RAPTOR" choice:"4" choice:"5" choice:"6" choice:"7" choice:"8" choice:"9" choice:"10" choice:"11" choice:"12"`
        //KeyPass string `short:"k" long:"KeyPass" description:"Your Key Passphase.  Recommend this be set in your archit configuration file" default:"insecure" env:"ARCHIT_KEYPASS"`
        //KeyPIN  int    `short:"n" long:"KeyPIN" description:"Your personal identificatio number, used to encrypt KeyPass" default:"0" env:"ARCHIT_KEYPIN"`

        //SeedMode        bool  `short:"S" long:"SeedMode" description:"Set Seed mode and bypass some checks & Wallet activity" enn:"ARCHIT_SEED"`
        //PortBase       int    `short:"B" long:"PortBase" description:"Primary port number Archit servers will listen to. Port# +1 will be used interally for server communication" default:"1958" env:"ARCHIT_PORT"`
        //WalletAddr     string `long:"WalletAddr" short:"A" description:"IMACredit Adddress for identity and transactions - !!!Do not use the default!!!" default:"9Ks***INVALID***Ady3pNn6jwbd9BT4Te" env:"ARCHIT_WALLETADDR"`
        //WalletIP       string `long:"WalletIP" short:"I" description:"IP name or address of IMACredit Wallet.  Recommend this be set in your archit configuration file" default:"localhost" env:"ARCHIT_WALLETIP"`
        //WalletPort     int    `long:"WalletPort" short:"P" description:"IMACredit Wallets's RPCPort setting." default:"64096" env:"ARCHIT_WALLETPORT"`
        //WalletUser     string `long:"WalletUser" short:"U" description:"IMACredit Wallet's RPCuser setting.  Recommend this be set in your archit configuration file" default:"ReplaceThis"`
        //WalletPassword string `long:"WalletPassword" short:"W" description:"IMACredit Wallet's RPCPassword settting.  HIGHLY recommend this be set in your archit configuration file" default:"AbsolutelyChangeThis"`

}

//  Allow them ALL, so user can play and see conf file settings.
var settingsCmd SettingsCommand

func init() {
	config.Parser.AddCommand("settings", "Displays current settings", 
		"Displays all current settings.  Use -S to show passwords", &settingsCmd)
}

func (ec *SettingsCommand) Execute(args []string) error {

	//util.Chaos = settingsCmd.Chaos
	//util.Raptor = settingsCmd.Raptor
	//util.KeyPass = settingsCmd.KeyPass
	//util.KeyPIN = settingsCmd.KeyPIN
	//util.SeedMode = settingsCmd.SeedMode
	//util.PortBase = settingsCmd.PortBase
	//util.WalletAddr = settingsCmd.WalletAddr
	//util.WalletIP = settingsCmd.WalletIP
	//util.WalletPort = settingsCmd.WalletPort
	//util.WalletUser = settingsCmd.WalletUser
	//util.WalletPassword = settingsCmd.WalletPassword
	config.Conf(false)

	util.Challenge()

	if !settingsCmd.NoHide {
		util.KeyPIN = -1
		util.KeyPass = "*Hidden*"
		util.WalletPassword = "*Hidden*"
	}
	log.Console("Settings potentially from the conf file, passed flags, or ENV variables:\n")
	log.Console("Core Settings:")
	log.Console("    Conf =",util.Conf)
	log.Console("    DBDir =",util.DBDir)
	log.Console("    LogFile =",util.LogFile)
	log.Console("    ResetLog =",util.ResetLog)
	log.Console("    Verbose =",util.Verbose)
	log.Console("\nOther SubCommand Settings from passed flags or ENV variables:")
	log.Console("    Chaos =",util.Chaos)
	log.Console("    KeyPass =",util.KeyPass)
	log.Console("    PortBase =",util.PortBase)
	log.Console("    Raptor =",util.Raptor)
	log.Console("    SeedMode =",util.SeedMode)
	log.Console("    WalletAddr =",util.WalletAddr)
	log.Console("    WalletIP =",util.WalletIP)
	log.Console("    WalletPassword =",util.WalletPassword)
	log.Console("    WalletPort =",util.WalletPort)
	log.Console("    WalletUser =",util.WalletUser)
	log.Console("\nValues ONLY storable in ENV Variables (and not recommended to do so!):")
	log.Console("    KeyPIN =",util.KeyPIN)

	return nil
}
