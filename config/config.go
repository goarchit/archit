// Process the configuration file, passed parameters, and ENV variables.
// Originally work created on 1/3/2017
//

package config

import (
	"fmt"
	"github.com/Unknwon/goconfig"
	"github.com/goarchit/archit/log"
	"github.com/goarchit/archit/util"
	"golang.org/x/crypto/scrypt"
	"net"
	"os"
	"strconv"
	"strings"
)

var clfirst bool 

func ParseCmdLine() {

	// Go parse the command line

	clfirst = true
	_, err := Parser.Parse()
	if err != nil {
		if strings.Contains(err.Error(), "Help Options:") {
			log.Debug("Help called.  Saw:", err)
			os.Exit(0)
		}
		if strings.Contains(err.Error(), "Please specify one command") {
			log.Debug("No comamnds or options specified.  Saw:", err)
			os.Exit(0)
		}
		fmt.Println("***Internal Parsing structure error: ",err,"***")
		os.Exit(1)
	}
}

func Conf(needKey bool) {

	if !clfirst {
		log.Critical("config.Conf called before config.ParsedCmdLine")
	}

	// And off to read the configuration file.  Command line parameters will need to trump

	// Need to patch up ~ in the configuration & log file name if it starts with that.
	util.Conf = util.FullPath(Archit.Conf)
	util.LogFile = util.FullPath(Archit.LogFile)

	conf, err := goconfig.LoadConfigFile(util.Conf)
	if err != nil {
		// Something is wrong, generate a good configuration file
		genGoodConfig()
		conf, err = goconfig.LoadConfigFile(util.Conf)
		if err != nil {
			log.Critical("Something is seriously wrong:",err)
		}
	}

	// Initialize log system 
	util.LogLevel = Archit.LogLevel
	util.ResetLog = Archit.ResetLog
	util.Verbose = Archit.Verbose
	log.Setup(util.LogLevel, util.LogFile, util.Verbose, util.ResetLog)

	// Start checking if the configuration file has a setting, and if a flag was set bu
	// default, override with configuration file.  If no entry is found in the configuration
	// file, and its a Core Setting, let it default
	// Deal with those Core Settings first (see output of "archit settings" for list

	value, err := conf.GetValue("", "DBDir")
	if err == nil {
		log.Trace("DBDir found in configuration file")
		o := Parser.FindOptionByLongName("DBDir")
		if o.IsSetDefault() {
			log.Trace("Configuration value of", value, "overriding default value", Archit.DBDir)
			Archit.DBDir = value
		}
	}
	util.DBDir = util.FullPath(Archit.DBDir)

	value, err = conf.GetValue("", "LogFile")
	if err == nil {
		log.Trace("Value of LogFile from config file:", value)
		o := Parser.FindOptionByLongName("LogFile")
		if o == nil {
			log.Critical("Internal Parser error finding LogFile")
		}
		if o.IsSetDefault() {
			log.Trace("Configuration value of", value, "overriding default value", util.LogFile)
			util.LogFile = value
			log.Setup(util.LogLevel, util.LogFile, util.Verbose, util.ResetLog)
		}
	}

	value, err = conf.GetValue("", "LogLevel")
	if err == nil {
		log.Trace("Value of LogLevel from config file:", value)
		o := Parser.FindOptionByLongName("LogLevel")
		if o == nil {
			log.Critical("Internal Parser error finding LogLevel")
		}
		if o.IsSetDefault() {
			log.Trace("Configuration value of", value, "overriding default value", util.LogLevel)
			util.LogLevel, err  = strconv.Atoi(value)
			if err != nil {
				log.Critical(err)
			}
			log.Setup(util.LogLevel, util.LogFile, util.Verbose, util.ResetLog)
		}
	}

	value, err = conf.GetValue("", "ResetLog")
	if err == nil {
		log.Trace("Value of ResetLog from config file:", value)
		o := Parser.FindOptionByLongName("ResetLog")
		if o == nil {
			log.Critical("Internal Parser error finding ResetLog")
		}
		if o.IsSetDefault() {
			log.Trace("Configuration value of", value, "overriding default value", util.ResetLog)
			rl, err := strconv.ParseBool(value)
			if err != nil {
				log.Critical("Error in conf file, ResetLog =",err)
			}
			util.ResetLog = rl
			log.Setup(util.LogLevel, util.LogFile, util.Verbose, util.ResetLog)
		}
	}

	value, err = conf.GetValue("", "Verbose")
	if err == nil {
		log.Trace("Value of Verbose from config file:", value)
		o := Parser.FindOptionByLongName("Verbose")
		if o == nil {
			log.Critical("Internal Parser error finding Verbose")
		}
		if o.IsSetDefault() {
			log.Trace("Configuration value of", value, "overriding default value", util.Verbose)
			util.Verbose, err  = strconv.Atoi(value)
			if err != nil {
				log.Critical(err)
			}
			if util.Verbose < 0 || util.Verbose > 2 {
				log.Critical("Error in conf file, Verbose must be between 0 and 2")
			}
			log.Setup(util.LogLevel, util.LogFile, util.Verbose, util.ResetLog)
		}
	}

	// Everything below MAY not have flags, depending on the command being issues
	// Note that NONE have a Archit. value, those only exist for Core Settings
	// Each command is responsibile to copy their parsed flags into the util. variable it
	// uses, so we are only concerned with conf file overrides

	value, err = conf.GetValue("", "PortBase")
	if err == nil {
		log.Trace("Value of Port from config file:", value)
		o := Parser.FindOptionByLongName("PortBase")
		if o == nil || o.IsSetDefault() {
			log.Trace("Configuration value of", value, "overriding default value", util.PortBase)
			util.PortBase, err = strconv.Atoi(value)
			if err != nil {
				log.Critical(err)
			}
		}
	}
	if util.SeedMode && (util.PortBase != 1958) {
		log.Critical("You want to be a seed?  PortBase must be set to 1958!")
	}

	// Go out and determine our public IP address
	util.PublicIP = net.JoinHostPort(util.GetExtIP(),strconv.Itoa(util.PortBase))
	log.Info("Farming at IP",util.PublicIP)

	//  Thats its... if you want to be a seed server...
	if util.SeedMode {
		log.Debug("SeedMode selected, exiting config.Conf() early")
		return
	}

	//  Otherwise proceed with checking everything else

	value, err = conf.GetValue("", "Raptor")
	if err == nil {
		log.Trace("Value of Raptor from config file:", value)
		o := Parser.FindOptionByLongName("Raptor")
		if o == nil || o.IsSetDefault() {
			log.Trace("Configuration value of", value, "overriding default value", util.Raptor)
			util.Raptor, err = strconv.Atoi(value)
			if err != nil {
				log.Critical(err)
			}
		}
	}
	util.Raptor = util.Raptor
	value, err = conf.GetValue("", "KeyPass")
	if err == nil {
		log.Trace("Value of KeyPass from config file:", value)
		o := Parser.FindOptionByLongName("KeyPass")
		if o == nil || o.IsSetDefault() {
			log.Trace("Configuration value of", value, "overriding default value", util.KeyPass)
			util.KeyPass = value
		}
	} else {
		log.Trace("KeyPass value NOT specified by config file")
		o := Parser.FindOptionByLongName("KeyPass")
		if o == nil || o.IsSetDefault() {
			log.Critical("KeyPass value must be set!")
		}
	}

	if needKey {
		// KeyPIN is a not allowed in the Configuration file
		// If a user REALLY wants to record it, they can do so as an ENV variable
		o := Parser.FindOptionByLongName("KeyPIN")
		if o == nil || o.IsSetDefault() {
			if util.KeyPIN == 0 {  //go-flag bug work aroundjj
			var ivalue int
			fmt.Printf("Default KeyPIN not allowed, please enter a KeyPIN: ")
			_, err := fmt.Scanf("%d", &ivalue)
			if err != nil {
				log.Critical("Error getting KeyPin:", err)
			}
			util.KeyPIN = ivalue
			}
		}
		// Generate and save DevivedKey
		util.DerivedKey, err = scrypt.Key([]byte(util.KeyPass), []byte(strconv.Itoa(util.KeyPIN)), 16384, 8, 1, 32)
		if err != nil {
			log.Critical("Error encrypting KeyPass:", err)
		}
		log.Trace("DerivedKey: ", util.DerivedKey, "len=", len(util.DerivedKey))
	}

	value, err = conf.GetValue("", "WalletAddr")
	if err == nil {
		log.Trace("Value of WalletAddr from config file:", value)
		o := Parser.FindOptionByLongName("WalletAddr")
		if o == nil || o.IsSetDefault() {
			log.Trace("Configuration value of", value, "overriding default value", util.WalletAddr)
			util.WalletAddr = value
		}
	} else {
		o := Parser.FindOptionByLongName("WalletAddr")
		if o == nil || o.IsSetDefault() {
			log.Console("Hey!  Don't use the default WalletAddr!!!")
			log.Console("Suggest adding 'WalletAddr =' in",util.Conf)
			log.Console("Issue an 'archit --help' and read the output for more hints")
			log.Critical("WalletAddr error")
		}
	}
	value, err = conf.GetValue("", "WalletIP")
	if err == nil {
		log.Trace("Value of WalletIP from config file:", value)
		o := Parser.FindOptionByLongName("WalletIP")
		if o == nil || o.IsSetDefault() {
			log.Trace("Configuration value of", value, "overriding default value", util.WalletIP)
			util.WalletIP = value
		}
	}
	value, err = conf.GetValue("", "WalletPort")
	if err == nil {
		log.Trace("Value of WalletPort from config file:", value)
		o := Parser.FindOptionByLongName("WalletPort")
		if o == nil || o.IsSetDefault() {
			log.Trace("Configuration value of", value, "overriding default value", util.WalletPort)
			util.WalletPort, err = strconv.Atoi(value)
			if err != nil {
				log.Critical(err)
			}
		}
	}
	value, err = conf.GetValue("", "WalletUser")
	if err == nil {
		log.Trace("Value of WalletUser from config file:", value)
		o := Parser.FindOptionByLongName("WalletUser")
		if o == nil || o.IsSetDefault() {
			log.Trace("Configuration value of", value, "overriding default value", util.WalletUser)
			util.WalletUser = value
		}
	} else {
		o := Parser.FindOptionByLongName("WalletUser")
		if o == nil || o.IsSetDefault() {
			log.Console("Hey!  Don't use the default WalletUser!!!")
			log.Console("Suggest adding 'WalletUser =' in",util.Conf)
			log.Console("Issue an 'archit --help' and read the output for more hints")
			log.Critical("WalletUser error")
		}
	}
	value, err = conf.GetValue("", "WalletPassword")
	if err == nil {
		log.Trace("WalletPassword found in configuration file")
		o := Parser.FindOptionByLongName("WalletPassword")
		if o == nil || o.IsSetDefault() {
			log.Trace("Configuration value of WalletPassword overriding default value")
			util.WalletPassword = value
		}
	} else {
		o := Parser.FindOptionByLongName("WalletPassword")
		if o == nil || o.IsSetDefault() {
			log.Console("Hey!  Don't use the default WalletPassword!!!")
			log.Console("Suggest adding 'WalletPassword =' in",util.Conf)
			log.Console("Issue an 'archit --help' and read the output for more hints")
			log.Critical("WalletPassword error")
		}
	}
	value, err = conf.GetValue("", "MinFreeSpace")
	if err == nil {
		log.Trace("Value of MinFreeSpace from config file:", value)
		o := Parser.FindOptionByLongName("MinFreeSpace")
		if o == nil || o.IsSetDefault() {
			log.Trace("Configuration value of", value, "overriding default value", util.MinFreeSpace)
			minFreeSpace, err := strconv.Atoi(value)
			if err != nil {
				log.Critical(err)
			}
			util.MinFreeSpace = uint64(minFreeSpace)
		}
	}
	value, err = conf.GetValue("", "DataDir")
	if err == nil {
		log.Trace("Value of DataDir from config file:", value)
		o := Parser.FindOptionByLongName("DataDir")
		if o == nil || o.IsSetDefault() {
			log.Trace("Configuration value of", value, "overriding default value", util.DataDir)
			util.DataDir = value
		}
	}
	util.DataDir = util.FullPath(util.DataDir)
	// Build the encryption Matrix
	util.BuildMatrix()
	// Make sure we have valud certification files
	private := util.DBDir + "/Private"
	public := util.DBDir + "/Public"
	// Check if cert files are available
	err = util.RSACheck(private, public)
	// If they are not, generate new ones
	if err != nil {
		log.Console("Generating security keys")
		util.RSAGenerate(private,public)
	}
	log.Trace("Wrapping up config.Conf(needKey), needKey:", needKey)
}
