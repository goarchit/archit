// Config command:  Displays your current configuration file
// Originally work created on 2/4/2017
//

package cmd

import (
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/log"
	"github.com/goarchit/archit/util"
	"io/ioutil"
	"fmt"
	"os"
)

type ConfigCommand struct {
}

func init() {
	configCmd := ConfigCommand{}
	_, err := config.Parser.AddCommand("config", "Displays your current configuration file", "", &configCmd)
        if err != nil {
                fmt.Println("Internal error parsing Config command:",err)
                os.Exit(1)
        }
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
