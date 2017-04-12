// Config command:  Displays your current configuration file
// Originally work created on 2/4/2017
//

package cmd

import (
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/log"
	"github.com/goarchit/archit/util"
	"io/ioutil"
)

type ConfigCommand struct {
}

func init() {
	configCmd := ConfigCommand{}
	config.Parser.AddCommand("config", "Displays your current configuration file", "", &configCmd)
}

func (ec *ConfigCommand) Execute(args []string) error {
	config.Conf(false)

	util.Challenge()

	log.Console(util.Conf+":\n")
	f, err := ioutil.ReadFile(util.Conf)
	if err != nil {
		log.Console("Unable to open specified configuration file.")
		return err
	} 
	if err != nil {
		log.Console("Error reading file:",err)
		return err
	}
	log.Console(string(f))
	return nil
}
