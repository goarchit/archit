// Settings command:  Display all configurable settings
// Originally work created on 1/8/2017
//

package cmd

import (
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/log"
)

type SettingsCommand struct{
}

func init() {
	settingsCmd := SettingsCommand{}
        config.Parser.AddCommand("settings","Displays all adjustable settings[Free]", "", &settingsCmd)
}

func (ec *SettingsCommand) Execute(args []string) error {
	config.Conf(false)
	log.Console("Settings command not implemented")
	return nil
}
