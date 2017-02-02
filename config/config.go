// Process the configuration file, passed parameters, and ENV variables.
// Originally work created on 1/3/2017
//

package config

import (
	"fmt"
	"github.com/Unknwon/goconfig"
	"github.com/goarchit/archit/log"
	"github.com/goarchit/archit/parser"
	"golang.org/x/crypto/scrypt"
	"os"
	"os/user"
	"strconv"
	"strings"
)

var DerivedKey string

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
			log.Debug("Configuration value of", value, "overriding default value", parser.Archit.Port)
			parser.Archit.Port, err = strconv.Atoi(value)
			if err != nil {
				log.Critical(err)
			}
		}
	}
	value, err = conf.GetValue("", "PublicIP")
	if err == nil {
		log.Debug("Value of PublicIP from config file:", value)
		o := parser.Parser.FindOptionByLongName("PublicIP")
		if o.IsSetDefault() {
			log.Debug("Configuration value of", value, "overriding default value", parser.Archit.PublicIP)
			parser.Archit.PublicIP = value
		}
	}
	value, err = conf.GetValue("", "LogLevel")
	if err == nil {
		log.Debug("Value of LogLevel from config file:", value)
		o := parser.Parser.FindOptionByLongName("LogLevel")
		if o.IsSetDefault() {
			log.Debug("Configuration value of", value, "overriding default value", parser.Archit.LogLevel)
			parser.Archit.Port, err = strconv.Atoi(value)
			if err != nil {
				log.Critical(err)
			}
		}
	}
	value, err = conf.GetValue("", "LogFile")
	if err == nil {
		log.Debug("Value of LogFIle from config file:", value)
		o := parser.Parser.FindOptionByLongName("LogFile")
		if o.IsSetDefault() {
			log.Debug("Configuration value of", value, "overriding default value", parser.Archit.LogFile)
			parser.Archit.LogFile = value
		}
	}
	value, err = conf.GetValue("", "KeyPass")
	if err == nil {
		log.Debug("Value of KeyPass from config file:", value)
		o := parser.Parser.FindOptionByLongName("KeyPass")
		if o.IsSetDefault() {
			log.Debug("Configuration value of", value, "overriding default value", parser.Archit.KeyPass)
			parser.Archit.KeyPass = value
		}
	} else {
		log.Debug("KeyPass value NOT specified by config file")
		o := parser.Parser.FindOptionByLongName("KeyPass")
		if o.IsSetDefault() {
			log.Critical("KeyPass value must be set!")
		}

	}
	value, err = conf.GetValue("", "KeyPIN")
	if err == nil {
		log.Debug("KeyPIN value specified by config file")
		o := parser.Parser.FindOptionByLongName("KeyPIN")
		if o.IsSetDefault() {
			log.Debug("Configuration value overriding default value")
			parser.Archit.KeyPIN, err = strconv.Atoi(value)
			if err != nil {
				log.Critical("Error parsing KeyPIN:", err)
			}
		}
	} else {
		log.Debug("KeyPIN value NOT specified by config file")
		o := parser.Parser.FindOptionByLongName("KeyPIN")
		if o.IsSetDefault() {
			var ivalue int
			log.Debug("KeyPIN value being asked of invoker")
			fmt.Printf("Please enter you KeyPIN: ")
			_, err := fmt.Scanf("%d", &ivalue)
			if err != nil {
				log.Critical("Error getting KeyPin:", err)
			}
			parser.Archit.KeyPIN = ivalue
		}

	}
	value, err = conf.GetValue("", "WalletAddr")
	if err == nil {
		log.Debug("Value of WalletAddr from config file:", value)
		o := parser.Parser.FindOptionByLongName("WalletAddr")
		if o.IsSetDefault() {
			log.Debug("Configuration value of", value, "overriding default value", parser.Archit.WalletAddr)
			parser.Archit.WalletAddr = value
		}
	}
	value, err = conf.GetValue("", "WalletIP")
	if err == nil {
		log.Debug("Value of WalletIP from config file:", value)
		o := parser.Parser.FindOptionByLongName("WalletIP")
		if o.IsSetDefault() {
			log.Debug("Configuration value of", value, "overriding default value", parser.Archit.WalletIP)
			parser.Archit.WalletIP = value
		}
	}
	value, err = conf.GetValue("", "WalletPort")
	if err == nil {
		log.Debug("Value of WalletPort from config file:", value)
		o := parser.Parser.FindOptionByLongName("WalletPort")
		if o.IsSetDefault() {
			log.Debug("Configuration value of", value, "overriding default value", parser.Archit.WalletPort)
			parser.Archit.WalletPort = value
		}
	}
	value, err = conf.GetValue("", "WalletUser")
	if err == nil {
		log.Debug("Value of WalletUser from config file:", value)
		o := parser.Parser.FindOptionByLongName("WalletUser")
		if o.IsSetDefault() {
			log.Debug("Configuration value of", value, "overriding default value", parser.Archit.WalletUser)
			parser.Archit.WalletUser = value
		}
	}
	value, err = conf.GetValue("", "WalletPassword")
	if err == nil {
		log.Debug("WalletPassword found in configuration file")
		o := parser.Parser.FindOptionByLongName("WalletPassword")
		if o.IsSetDefault() {
			log.Debug("Configuration value of WalletPassword overriding default value")
			parser.Archit.WalletPassword = value
		}
	}
	// Generate and save DevivedKey
	DerivedKey, err := scrypt.Key([]byte(parser.Archit.KeyPass), []byte(strconv.Itoa(parser.Archit.KeyPIN)), 32768, 16, 2, 32)
	if err != nil {
		log.Critical("Error encrypting KeyPass:", err)
	}
	log.Debug("DerivedKey: ", DerivedKey)
	log.Debug("Wrapping up config.Read()")
}
