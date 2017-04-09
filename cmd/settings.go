// Settings command:  Display all configurable settings
// Originally work created on 1/8/2017
//

package cmd

import (
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/log"
	"github.com/goarchit/archit/util"
	"encoding/json"
)

type SettingsCommand struct {
	NoHide	bool `short:"S" long:"ShowPasswords" description:"Show passwords" env:"ARCHIT_NOHIDE"`
}

var settingsCmd SettingsCommand

func init() {
	config.Parser.AddCommand("settings", "Displays current settings", 
		"Displays all current settings.  Use -S to show passwords", &settingsCmd)
}

func (ec *SettingsCommand) Execute(args []string) error {
	config.Conf(false)

	util.Challenge()

	if !settingsCmd.NoHide {
		config.Archit.KeyPass = "*Hidden*"
		config.Archit.WalletPassword = "*Hidden*"
	}
	s, err := json.MarshalIndent(config.Archit,"","   ")
	if err != nil {
		log.Critical("Internal issue:",err)
	}
	log.Console("Parsed values:",string(s))
	return nil
}
