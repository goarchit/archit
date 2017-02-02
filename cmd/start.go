// Start command:  e.g.  archit start
// Originally work created on 1/8/2017
//

package cmd

import (
	"github.com/goarchit/archit/log"
)

type StartCommand struct{
}

func (ec *StartCommand) Execute(args []string) error {
	log.Debug("In StartCommand")
	return nil
}
