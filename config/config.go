// Process the configuration file, passed parameters, and ENV variables.
// Originally work created on 1/3/2017
//

package config

import (
	"fmt"
	"github.com/Unknwon/goconfig"
	"github.com/goarchit/archit/log"
	"golang.org/x/crypto/scrypt"
	"os"
	"os/user"
	"strconv"
	"strings"
)

var DerivedKey string
var clfirst bool = false

func ParseCmdLine() {

	// Go parse the command line

	clfirst = true
	_, err := Parser.Parse()
	if err != nil {
		if strings.Contains(err.Error(), "Help Options:") {
			log.Debug("Help called.  Saw:", err)
			os.Exit(0)
		}
		os.Exit(1)
	}
}

func Conf(needKey bool) {

	var verbose int = 0

	if !clfirst {
		log.Critical("config.Conf called before config.ParsedCmdLine")
	}

	// And off to read the configuration file.  Command line parameters will need to trump

	// Need to patch up ~ in the configuration file name if it starts with that.
	if len(Archit.Conf) >= 2 && Archit.Conf[:2] == "~/" {
		usr, _ := user.Current()
		dir := usr.HomeDir + "/"
		Archit.Conf = strings.Replace(Archit.Conf, "~/", dir, 1)
	}
	conf, err := goconfig.LoadConfigFile(Archit.Conf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Initialize log system and write Informational messages
	if Archit.Verbose {
		verbose = 1
	}
	if Archit.VVerbose {
		verbose = 2
	}
	log.Setup(Archit.LogLevel, Archit.LogFile, verbose, Archit.ResetLog)

	// Start checking if the configuration file has a setting, and if a f;ag was set bu
	// default, override with configuration file
	value, err := conf.GetValue("", "PortBase")
	if err == nil {
		log.Debug("Value of Port from config file:", value)
		o := Parser.FindOptionByLongName("PortBase")
		if o.IsSetDefault() {
			log.Debug("Configuration value of", value, "overriding default value", Archit.PortBase)
			Archit.PortBase, err = strconv.Atoi(value)
			if err != nil {
				log.Critical(err)
			}
		}
	}
	value, err = conf.GetValue("", "PublicIP")
	if err == nil {
		log.Debug("Value of PublicIP from config file:", value)
		o := Parser.FindOptionByLongName("PublicIP")
		if o.IsSetDefault() {
			log.Debug("Configuration value of", value, "overriding default value", Archit.PublicIP)
			Archit.PublicIP = value
		}
	}
	value, err = conf.GetValue("", "LogLevel")
	if err == nil {
		log.Debug("Value of LogLevel from config file:", value)
		o := Parser.FindOptionByLongName("LogLevel")
		if o.IsSetDefault() {
			log.Debug("Configuration value of", value, "overriding default value", Archit.LogLevel)
			Archit.LogLevel, err = strconv.Atoi(value)
			if err != nil {
				log.Critical(err)
			}
		}
	}
	value, err = conf.GetValue("", "LogFile")
	if err == nil {
		log.Debug("Value of LogFIle from config file:", value)
		o := Parser.FindOptionByLongName("LogFile")
		if o.IsSetDefault() {
			log.Debug("Configuration value of", value, "overriding default value", Archit.LogFile)
			Archit.LogFile = value
		}
	}
	if needKey {
		value, err = conf.GetValue("", "KeyPass")
		if err == nil {
			log.Debug("Value of KeyPass from config file:", value)
			o := Parser.FindOptionByLongName("KeyPass")
			if o.IsSetDefault() {
				log.Debug("Configuration value of", value, "overriding default value", Archit.KeyPass)
				Archit.KeyPass = value
			}
		} else {
			log.Debug("KeyPass value NOT specified by config file")
			o := Parser.FindOptionByLongName("KeyPass")
			if o.IsSetDefault() {
				log.Critical("KeyPass value must be set!")
			}

		}
		value, err = conf.GetValue("", "KeyPIN")
		if err == nil {
			log.Debug("KeyPIN value specified by config file")
			o := Parser.FindOptionByLongName("KeyPIN")
			if o.IsSetDefault() {
				log.Debug("Configuration value overriding default value")
				Archit.KeyPIN, err = strconv.Atoi(value)
				if err != nil {
					log.Critical("Error parsing KeyPIN:", err)
				}
			}
		} else {
			log.Debug("KeyPIN value NOT specified by config file")
			o := Parser.FindOptionByLongName("KeyPIN")
			if o.IsSetDefault() {
				var ivalue int
				log.Debug("KeyPIN value being asked of invoker")
				fmt.Printf("Please enter you KeyPIN: ")
				_, err := fmt.Scanf("%d", &ivalue)
				if err != nil {
					log.Critical("Error getting KeyPin:", err)
				}
				Archit.KeyPIN = ivalue
			}

		}
		// Generate and save DevivedKey
		DerivedKey, err := scrypt.Key([]byte(Archit.KeyPass), []byte(strconv.Itoa(Archit.KeyPIN)), 32768, 16, 2, 32)
		if err != nil {
			log.Critical("Error encrypting KeyPass:", err)
		}
		log.Debug("DerivedKey: ", DerivedKey)
	}
	value, err = conf.GetValue("", "WalletAddr")
	if err == nil {
		log.Debug("Value of WalletAddr from config file:", value)
		o := Parser.FindOptionByLongName("WalletAddr")
		if o.IsSetDefault() {
			log.Debug("Configuration value of", value, "overriding default value", Archit.WalletAddr)
			Archit.WalletAddr = value
		}
	}
	value, err = conf.GetValue("", "WalletIP")
	if err == nil {
		log.Debug("Value of WalletIP from config file:", value)
		o := Parser.FindOptionByLongName("WalletIP")
		if o.IsSetDefault() {
			log.Debug("Configuration value of", value, "overriding default value", Archit.WalletIP)
			Archit.WalletIP = value
		}
	}
	value, err = conf.GetValue("", "WalletPort")
	if err == nil {
		log.Debug("Value of WalletPort from config file:", value)
		o := Parser.FindOptionByLongName("WalletPort")
		if o.IsSetDefault() {
			log.Debug("Configuration value of", value, "overriding default value", Archit.WalletPort)
			Archit.WalletPort, err = strconv.Atoi(value)
			if err != nil {
				log.Critical(err)
			}
		}
	}
	value, err = conf.GetValue("", "WalletUser")
	if err == nil {
		log.Debug("Value of WalletUser from config file:", value)
		o := Parser.FindOptionByLongName("WalletUser")
		if o.IsSetDefault() {
			log.Debug("Configuration value of", value, "overriding default value", Archit.WalletUser)
			Archit.WalletUser = value
		}
	}
	value, err = conf.GetValue("", "WalletPassword")
	if err == nil {
		log.Debug("WalletPassword found in configuration file")
		o := Parser.FindOptionByLongName("WalletPassword")
		if o.IsSetDefault() {
			log.Debug("Configuration value of WalletPassword overriding default value")
			Archit.WalletPassword = value
		}
	}
	value, err = conf.GetValue("", "DBDir")
	if err == nil {
		log.Debug("DBDir found in configuration file")
		o := Parser.FindOptionByLongName("DBDir")
		if o.IsSetDefault() {
			log.Debug("Configuration value of", value, "overriding default value", Archit.DBDir)
			Archit.DBDir = value
		}
	}
	// Need to patch up ~ in the database file name if it starts with that.
	if len(Archit.DBDir) >= 2 && Archit.DBDir[:2] == "~/" {
		usr, _ := user.Current()
		dir := usr.HomeDir + "/"
		Archit.DBDir = strings.Replace(Archit.DBDir, "~/", dir, 1)
	}
	log.Debug("Wrapping up config.Conf(needKey), needKey:", needKey)
}
