// Start command:  e.g.  archit start
// Originally work created on 1/8/2017
//

package cmd

import (
	"github.com/goarchit/archit/parser"
)

func Define() {
	startCmd := StartCommand{}
	parser.Parser.AddCommand("start","Starts the ArchIt Network server", "", &startCmd)
	statusCmd := StatusCommand{}
	parser.Parser.AddCommand("status","Shows the status of the ArchIt Network server", "", &statusCmd)
}
