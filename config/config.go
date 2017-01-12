// Process the configuration file, passed parameters, and ENV variables.
// Originally work created on 1/3/2017
//

package config

import (
	"fmt"
	"github.com/Unknwon/goconfig"
	"github.com/goarchit/archit/log"
	"github.com/goarchit/archit/parser"
	"os"
	"os/user"
	"strings"
	"strconv"
)

func Read() {

	var verbose int = 0

	// Go parse the command line

	_, err := parser.Parser.Parse()
	if err != nil {
		if strings.Contains(err.Error(), "Help Options:") {
			log.Debug("Help called.  Saw:", err)
			os.Exit(0)
		}
		log.Critical(err)
	}
	// And off to read the configuration file.  Command line parameters will need to trump

	// Need to patch up ~ in the configuration file name if it starts with that.
	if len(parser.Archit.Conf) >= 2 && parser.Archit.Conf[:2] == "~/" {
		usr, _ := user.Current()
		dir := usr.HomeDir + "/"
		parser.Archit.Conf = strings.Replace(parser.Archit.Conf, "~/", dir, 1)
	}
	conf, err := goconfig.LoadConfigFile(parser.Archit.Conf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Initialize log system and write Informational messages
	if parser.Archit.Verbose {
		verbose = 1
	}
	if parser.Archit.VVerbose {
		verbose = 2
	}
	log.Setup(parser.Archit.LogLevel, parser.Archit.LogFile, verbose, parser.Archit.ResetLog)

	// Start checking if the configuration file has a setting, and if a f;ag was set bu
	// default, override with configuration file
	value, err := conf.GetValue("", "Port")
	if err == nil {
		log.Debug("Value of Port from config file:", value)
		o := parser.Parser.FindOptionByLongName("Port")
		if o.IsSetDefault() {
			log.Debug("Configure value of", value, "overriding default value", parser.Archit.Port)
			parser.Archit.Port, err = strconv.Atoi(value)
			if err != nil{
				log.Critical(err)
			}
		}
	}
	value, err = conf.GetValue("", "PublicIP")
	if err == nil {
		log.Debug("Value of PublicIP from config file:", value)
		o := parser.Parser.FindOptionByLongName("PublicIP")
		if o.IsSetDefault() {
			log.Debug("Configure value of", value, "overriding default value", parser.Archit.PublicIP)
			parser.Archit.PublicIP = value
		}
	}
	value, err = conf.GetValue("", "LogLevel")
	if err == nil {
		log.Debug("Value of LogLevel from config file:", value)
		o := parser.Parser.FindOptionByLongName("LogLevel")
		if o.IsSetDefault() {
			log.Debug("Configure value of", value, "overriding default value", parser.Archit.LogLevel)
			parser.Archit.Port, err = strconv.Atoi(value)
			if err != nil{
				log.Critical(err)
			}
		}
	}
	value, err = conf.GetValue("", "LogFile")
	if err == nil {
		log.Debug("Value of LogFIle from config file:", value)
		o := parser.Parser.FindOptionByLongName("LogFile")
		if o.IsSetDefault() {
			log.Debug("Configure value of", value, "overriding default value", parser.Archit.LogFile)
			parser.Archit.LogFile = value
		}
	}
	log.Debug("Wrapping up config.Read()")
}
