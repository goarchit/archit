// Cron.go - deal with all things that should happen occassionally
package db

import (
	"github.com/goarchit/archit/farmer"
	"github.com/goarchit/archit/log"
	"github.com/goarchit/archit/util"
	"net"
	"time"
)

func CronHourly() {
	log.Trace("Internal CronHourly() called")
	// loop forever //
	for {
		t := time.NewTimer(1 * time.Hour)
		<-t.C
		//	db.FlushPeerMap()	//  Update the PeerMap in the bolt database
	}
}

func CronDaily() {
	log.Trace("Internal CronDaily() called")
	// loop forever //
	for {
		t := time.NewTimer(24 * time.Hour)
		<-t.C
		oldHost, _, err := net.SplitHostPort(util.ServerIP)
		if err != nil {
			log.Critical("Checkip:  Error spliting host & port -", err)
		}
		newHost := util.GetExtIP()
		if oldHost != newHost {
			log.Warning("Public IP address has changed!!!")
			log.Warning("Old IP:", oldHost, "New IP:", newHost)
			log.Warning("Attemptint to stop farmer process")
			util.FarmerStop <- true
			log.Warning("Sleeping 10 seconds before attempting restart")
			time.Sleep(10 * time.Second)
			go farmer.Run(util.FarmerStop)
			log.Warning("Restart initiated, good luck!")
		}
	}
}
