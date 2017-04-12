// ArchIt main routine
// Originally work created on 1/3/2017
//

package config

import (
	"github.com/jessevdk/go-flags"
)

//  Global command line flags
type architCommand struct {
	Conf           string `short:"c" long:"Conf" description:"Name of the conf file" default:"~/.archit/archit.conf" env:"ARCHIT_CONF"`
	DBDir          string `short:"b" long:"DBDir" description:"Path for persistance databases" default:"~/.archit" env:"ARCHIT_DBDIR"`
	LogFile        string `short:"l" long:"LogFile" description:"Name of your logfile" default:"archit.log" env:"ARCHIT_LOGFILE"`
	LogLevel       int    `short:"d" long:"LogLevel" description:"Logging level" default:"3" choice:"0" choice:"1" choice:"2" choice:"3" choice:"4" choice:"5" env:"ARCHIT_LOGLEVEL"`
	ResetLog       bool   `short:"r" long:"ResetLog" description:"Reset the logfile when starting" env:"ARCHIT_RESETLOG"`
	Verbose        int   `short:"v" long:"Verbose" description:"Show additional messages" choice:"0" choice:"1" choice:"2" default:"0" env:"ARCHIT_VERBOSE"`
}

var Archit architCommand

var Parser *flags.Parser = flags.NewParser(&Archit, flags.Default)
