// Audit command - Performs various Archit Network audits
// Originally work created on 2/4/2017
//

package cmd

import (
	"github.com/ricochet2200/go-disk-usage/du"
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/log"
	"github.com/goarchit/archit/util"
	"fmt"
	"os"
	"net"
)

type TestCommand struct {
}

func init() {
	testCmd := TestCommand{}
	_,err := config.Parser.AddCommand("test", "Go functionality testing", "", &testCmd)
        if err != nil {
                fmt.Println("Internal error parsing Test command:",err)
                os.Exit(1)
        }
}

var KB = uint64(1024)

func (ec *TestCommand) Execute(args []string) error {
	config.Conf(false)
	rifs := util.RoutedInterface("ip", net.FlagUp|net.FlagBroadcast)
	if rifs != nil {
		log.Console("Routed interface is ", rifs.HardwareAddr.String())
	}
	log.Console("Internal IP: ",util.GetOutboundIP())
	log.Console("External IP: ",util.GetExtIP())

	usage := du.NewDiskUsage(".")
	fmt.Println("Free:", usage.Free()/(KB*KB))
	fmt.Println("Available:", usage.Available()/(KB*KB))
	fmt.Println("Size:", usage.Size()/(KB*KB))
	fmt.Println("Used:", usage.Used()/(KB*KB))
	fmt.Println("Usage:", usage.Usage()*100, "%")

	log.Console("Testing complete")

	return nil
}
