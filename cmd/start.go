// Start command:  e.g.  archit start
// Originally work created on 1/8/2017
//

package cmd

import (
	"fmt"
)

type StartCommand struct{
}

func (ec *StartCommand) Execute(args []string) error {
	fmt.Println("In StartCommand")
	return nil
}
