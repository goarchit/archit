// Config command:  Displays your current configuration file
// Originally work created on 2/4/2017
//

package cmd

import (
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/log"
)

type ConfigCommand struct {
}

func init() {
	configCmd := ConfigCommand{}
	config.Parser.AddCommand("config", "Displays your current configuration file[Free]", "", &configCmd)
}

func (ec *ConfigCommand) Execute(args []string) error {
	config.Conf(false)
	log.Console("Config command not implemented")
	return nil
}
