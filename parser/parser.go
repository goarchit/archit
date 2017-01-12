// ArchIt main routine
// Originally work created on 1/3/2017
//

package parser

import (
	"github.com/jessevdk/go-flags"
)

type ArchitCommand struct {
	Port     int    `short:"p" long:"Port" description:"Primary port number Archit server will listen to. Port# +1 will be used interally for server communication" default:"1958" env:"ARCHIT_PORT"`
	PublicIP string `short:"i" long:"PublicIP" description:"IP address for server to use (typically your public IP)" default:"localhost" env:"ARCHIT_PUBLICIP"`
	LogLevel int    `short:"d" long:"LogLevel" description:"Message level (0-4)" default:"3" choice:"0" choice:"1" choice:"2" choice:"3" choice:"4" env:"ARCHIT_LOGLEVEL"`
	LogFile  string `short:"l" long:"LogFile" description:"Name of your logfile" default:"archit.log" env:"ARCHIT_LOGFILE"`
	ResetLog bool `short:"r" long:"ResetLog" description:"Reset the logfile instead of appending to it" env:"ARCHIT_RESETLOG"`
	Conf     string `short:"c" long:"Conf" description:"Configuration file name" default:"~/.archit/archit.conf" env:"ARCHIT_CONF"`
	Verbose bool `short:"v" long:"Verbose" description:"Show Informational messages on the console" env:"ARCHIT_VERBOSE"`
	VVerbose bool `short:"V" long:"VeryVerbose" description:"Show ALL messages, include debug if set" env:"ARCHIT_DEBUG"`
}

var Archit ArchitCommand

var Parser *flags.Parser = flags.NewParser(&Archit, flags.Default)
