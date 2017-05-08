// Cron.go - deal with all things that should happen occassionally
package farmer

import (
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
		FlushPeerMap()	//  Update the PeerMap in the bolt database
		CheckPeers()	//  Adjust reputation of peers that don't Ping
	}
}

func CronDaily() {
	log.Trace("Internal CronDaily() called")
	// loop forever //
	for {
		t := time.NewTimer(24 * time.Hour)
		<-t.C
		oldHost, _, err := net.SplitHostPort(util.PublicIP)
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
			go Run(util.FarmerStop)
			log.Warning("Restart initiated, good luck!")
		}
	}
}
