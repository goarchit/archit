// ArchIt logging routine
// Originally work created on 1/7/2017
//

package server

import (
	"github.com/boltdb/bolt"
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/log"
	"fmt"
	"time"
)

var PeersDB *bolt.DB
var TrackingDB *bolt.DB
var err error

func DB(c chan string) {
	// Start the DBs
	PeersDB, err = bolt.Open(config.Archit.DBDir+"/Peers.bolt", 0600, &bolt.Options{Timeout: 2 * time.Second})
	if err != nil {
		log.Critical("DB Error opening Peers.bolt:", err)
	}
	defer PeersDB.Close()
	TrackingDB, err = bolt.Open(config.Archit.DBDir+"/Tracking.bolt", 0600, &bolt.Options{Timeout: 2 * time.Second})
	if err != nil {
		log.Critical("DB Error opening Tracking.bolt:", err)
	}
	defer TrackingDB.Close()

	// Now go process commands
	cmd := <-c
	for cmd != "stop" {
		switch cmd {
		case "status":
			c <- dbStatus()
		default:
			log.Error("DB passed unknown command", cmd)
		}
		cmd = <-c
	}
	c <- "DB shutdown complete"
}

func dbStatus() string {
	stats := PeersDB.Stats()
	response := fmt.Sprintf("Peers.bolt Reads: %d, Writes: %d, Nodes: %d, Rebalances: %d  Splits: %d\n",
		stats.TxN, stats.TxStats.Write, stats.TxStats.NodeCount, stats.TxStats.Rebalance, stats.TxStats.Split)
	stats = TrackingDB.Stats()
	response += fmt.Sprintf("Tracking.bolt Reads: %d, Writes: %d, Nodes: %d, Rebalances: %d  Splits: %d\n",
		stats.TxN, stats.TxStats.Write, stats.TxStats.NodeCount, stats.TxStats.Rebalance, stats.TxStats.Split)
	return response
}
