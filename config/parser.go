// ArchIt main routine
// Originally work created on 1/3/2017
//

package config

import (
	"github.com/jessevdk/go-flags"
)

type architCommand struct {
	Conf           string `short:"c" long:"Conf" description:"Name of the conf file" default:"~/.archit/archit.conf" env:"ARCHIT_CONF"`
	DBDir          string `short:"b" long:"DBDir" description:"Path for persistance databases" default:"~/.archit" env:"ARCHIT_DBDIR"`
	LogFile        string `short:"l" long:"LogFile" description:"Name of your logfile" default:"archit.log" env:"ARCHIT_LOGFILE"`
	LogLevel       int    `short:"d" long:"LogLevel" description:"Logging level" default:"3" choice:"0" choice:"1" choice:"2" choice:"3" choice:"4" choice:"5" env:"ARCHIT_LOGLEVEL"`
	ResetLog       bool   `short:"r" long:"ResetLog" description:"Reset the logfile when starting" env:"ARCHIT_RESETLOG"`
	Verbose        int   `short:"v" long:"Verbosity" description:"Show additional messages" choice:"0" choice:"1" choice:"2" default:"0" env:"ARCHIT_VERBOSE"`
	// Hidden Values
        Raptor  int    `short:"R" long:"Raptor" default:"8" env:"ARCHIT_RAPTOR" choice:"4" choice:"5" choice:"6" choice:"7" choice:"8" choice:"9" choice:"10" choice:"11" choice:"12" hidden:"Y"`
	KeyPass string `short:"K" long:"KeyPass" default:"insecure" env:"ARCHIT_KEYPASS" hidden:"Y"`
	KeyPIN  int    `short:"N" long:"KeyPIN"  default:"0" env:"ARCHIT_KEYPIN" hidden:"Y"`
	PortBase       int    `short:"B" long:"PortBase" default:"1958" env:"ARCHIT_PORT" hidden:"Y"`
	WalletAddr     string `short:"A" long:"WalletAddr" default:"9KsqKMgLjzBWCidAAdy3pNn6jwbd9BT4Te" env:"ARCHIT_WALLETADDR" hidden:"Y"`
	WalletIP       string `short:"I" long:"WalletIP" default:"localhost" env:"ARCHIT_WALLETIP" hidden:"Y"`
	WalletPort     int    `short:"P" long:"WalletPort" default:"64096" env:"ARCHIT_WALLETPORT" hidden:"Y"`
	WalletUser     string `short:"U" long:"WalletUser" default:"ReplaceThis" hidden:"Y"`
	WalletPassword string `short:"W" long:"WalletPassword" default:"AbsolutelyChangeThis" hidden:"Y"`
}

var Archit architCommand

var Parser *flags.Parser = flags.NewParser(&Archit, flags.Default)
