// Audit command - Performs various Archit Network audits
// Originally work created on 2/4/2017
//

package cmd

import (
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/log"
	"fmt"
	"os"
)

type AuditCommand struct {
}

func init() {
	auditCmd := AuditCommand{}
	_,err := config.Parser.AddCommand("audit", "Perform various Archit network auditing tasks[Fee!]", "", &auditCmd)
	if err != nil {
                fmt.Println("Internal error parsing Audit command:",err)
                os.Exit(1)
        }
}

func (ec *AuditCommand) Execute(args []string) error {
	config.Conf(true)
	log.Console("Auditing is not yet fully implemented")
	return nil
}
