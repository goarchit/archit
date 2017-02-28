// Audit command - Performs various Archit Network audits
// Originally work created on 2/4/2017
//

package cmd

import (
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/log"
	"github.com/goarchit/archit/util"
	"net"
	"os"
)

type TestCommand struct {
}

func init() {
	testCmd := TestCommand{}
	config.Parser.AddCommand("test", "Go functionality testing", "", &testCmd)
}

func (ec *TestCommand) Execute(args []string) error {
	config.Conf(false)
	rifs := util.RoutedInterface("ip", net.FlagUp|net.FlagBroadcast)
	if rifs != nil {
		log.Console("Routed interface is ", rifs.HardwareAddr.String())
	}

	log.Console("Starting encryption test")
	f, err := os.Open("LICENSE.md")
	if err != nil {
		log.Critical(err)
	}
	buffer := make([]byte, 32*1024)
	n1, err := f.Read(buffer)
	if err != nil {
		log.Critical(err)
	}
	if n1 != 32*1024 {
		log.Critical("LICENSE.md only contains",n1,"bytes")
	}
	tbuffer := make([]byte, 40*1024) 
	copy(tbuffer,buffer)  // Initialize the first 32 blocks
	copy(tbuffer[32*1024:],buffer[:8*1024-1])
	for i:=0; i<32; i++ {
		go encode(&tbuffer, &buffer, i)
	}

	log.Console("Buffer len",len(buffer),"tbuffer len",len(tbuffer))
	log.Console(string(tbuffer))
	log.Console("Testing complete")

	return nil
}

func encode(tblock, block *[]byte, i int) {
	tstart := i*1024
	tend := tstart + 1023
	for j:=1; j<7; j++ {
		bstart := (i+j)*1024;
		bend := bstart + 1023
		util.XorBytes(&tblock[tstart:tend],&tblock[tstart:tend],&block[bstart:bend])
	}
}
