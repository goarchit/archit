// Start command:  e.g.  archit start
// Originally work created on 1/8/2017
//

package cmd

import (
	"fmt"
)

type StatusCommand struct{
}

func (ec *StatusCommand) Execute(args []string) error {
	fmt.Println("In StatusCommand")
	return nil
}
