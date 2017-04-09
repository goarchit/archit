// Start command:  Start internal services
// Originally work created on 1/8/2017
//

package cmd

import (
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/log"
	"github.com/goarchit/archit/farmer"
	"github.com/goarchit/archit/util"
	"os"
	"os/signal"
	"io/ioutil"
	"strconv"
)

type FarmCommand struct{
}

func init() {
	farmCmd := FarmCommand{}
        config.Parser.AddCommand("farm","Starts the ArchIt farming & wallet servers[Free]", "", &farmCmd)
}

func (ec *FarmCommand) Execute(args []string) error {
	config.Conf(false)
	pid := os.Getpid()
	err := ioutil.WriteFile(config.Archit.DBDir+"/archit.pid", []byte(strconv.Itoa(pid)), 0644)
	if err != nil {
		log.Critical(err)
	}
	log.Console("Starting Farmer node... Ctrl-C to stop or kill pid", pid)
	go farmer.Run(util.FarmerStop)   //  Loops forever waiting on incomming web request
	//  Wait for a Ctrl-C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c  
	log.Console("\nTrying a clean shutdown")
	util.FarmerStop <- true
	return nil
}
