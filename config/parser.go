// ArchIt main routine
// Originally work created on 1/3/2017
//

package config

import (
	"github.com/jessevdk/go-flags"
)

type architCommand struct {
	Conf           string `short:"c" long:"Conf" description:"Location and name of the configuration file" default:"~/.archit/archit.conf" env:"ARCHIT_CONF"`
	DBDir     string `short:"b" long:"DBDir" description:"Location for persistance databases" default:"~/.archit" env:"ARCHIT_DBDIR"`
	LogFile        string `short:"l" long:"LogFile" description:"Name of your logfile" default:"archit.log" env:"ARCHIT_LOGFILE"`
	LogLevel       int    `short:"d" long:"LogLevel" description:"Message level (0-5)" default:"3" choice:"0" choice:"1" choice:"2" choice:"3" choice:"4" choice:"5" env:"ARCHIT_LOGLEVEL"`
	KeyPass        string `short:"k" long:"KeyPass" description:"Your Key Passphase.  Recommend this be set in your archit configuration file" default:"insecure" env:"ARCHIT_KEYPASS"`
	KeyPIN         int    `short:"n" long:"KeyPIN" description:"Your personal identificatio number, used to encrypt KeyPass" default:"0" env:"ARCHIT_KEYPIN"`
	PortBase       int    `short:"p" long:"PortBase" description:"Primary port number Archit servers will listen to. Port# +1 will be used interally for server communication" default:"1958" env:"ARCHIT_PORT"`
	Raptor		int	`short:"R" long:"Raptor" description:"Raptor factor - how many extra fountain blocks to generate (4-12)" default:"8" env:"ARCHIT_RAPTOR" choice:"4" choice:"5" choice:"6" choice:"7" choice:"8" choice:"9" choice:"10" choice:"11" choice:"12"`
	ResetLog       bool   `short:"r" long:"ResetLog" description:"Reset the logfile instead of appending to it" env:"ARCHIT_RESETLOG"`
	Verbose        bool   `short:"v" long:"Verbose" description:"Show Informational messages on the console" env:"ARCHIT_VERBOSE"`
	VVerbose       bool   `short:"V" long:"VeryVerbose" description:"Show ALL messages, include debug if set" env:"ARCHIT_DEBUG"`
	WalletAddr     string `long:"WalletAddr" short:"A" description:"IMACredit Adddress for identity and transactions - !!!Do not use the default!!!" default:"9KsqKMgLjzBWCidAAdy3pNn6jwbd9BT4Te" env:"ARCHIT_WALLETADDR"`
	WalletIP       string `long:"WalletIP" short:"I" description:"IP name or address of IMACredit Wallet.  Recommend this be set in your archit configuration file" default:"localhost" env:"ARCHIT_WALLETIP"`
	WalletPort     int `long:"WalletPort" short:"P" description:"IMACredit Wallets's RPCPort setting." default:"64096" env:"ARCHIT_WALLETPORT"`
	WalletUser     string `long:"WalletUser" short:"U" description:"IMACredit Wallet's RPCuser setting.  Recommend this be set in your archit configuration file" default:"ReplaceThis"`
	WalletPassword string `long:"WalletPassword" short:"W" description:"IMACredit Wallet's RPCPassword settting.  HIGHLY recommend this be set in your archit configuration file" default:"AbsolutelyChangeThis"`
	Chaos		bool `long:"Chaos" short:"C" description:"Add a check for a random bit flip (chaos bit) after encoding and decoding - this doubles the computational load for sending data and normally is not needed" env:"ARCHIT_CHAOS"`
}

var Archit architCommand

var Parser *flags.Parser = flags.NewParser(&Archit, flags.Default)
