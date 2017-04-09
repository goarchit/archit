// DB command:  Display local Archit.bolt stats
// Originally work created on 3/19/2017
//

package cmd

import (
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/db"
	"github.com/goarchit/archit/log"
)

type DbStatusCommand struct {
}

func init() {
	dbstatusCmd := DbStatusCommand{}
	config.Parser.AddCommand("dbstatus", "Shows the status of the renters database", "", &dbstatusCmd)
}

func (ec *DbStatusCommand) Execute(args []string) error {
	config.Conf(false)
	db.Open()
	defer db.Close()
	response := db.Status()
	log.Console(response)
	return nil
}
