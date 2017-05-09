// ArchIt Farming DB routines
// Code evolved along with client code
//
// These routines handle everything associated with peer information.  A client instance
// will request data via RPC to access that information

// Note the ../db/db.go routine handles all fileinfo data

package farmer

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/goarchit/archit/log"
	"github.com/goarchit/archit/util"
	"path/filepath"
	"time"
)

const PeersBucket = "PeersBucket"

var peersDB *bolt.DB
var err error

func DB(c chan string) {
	var peerbytes bytes.Buffer

	// Start the DBs
	dbName := filepath.Join(util.DBDir, util.PeerDBName)
	peersDB, err = bolt.Open(dbName, 0600, &bolt.Options{Timeout: 2 * time.Second})
	if err != nil {
		log.Critical("DB Error opening PeerInfo.bolt:", err)
	}
	defer peersDB.Close()
	// Handle buckets
	err = peersDB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(PeersBucket))
		if err != nil {
			log.Critical("Error creating bucket:", PeersBucket)
		}
		// Load up the reputation matrix
		peerbytes.Write(b.Get([]byte("PeerMap")))
		return nil
	})
	if peerbytes.Len() != 0 {
		dec := gob.NewDecoder(&peerbytes)
		err := dec.Decode(&PeerMap.PL)
		if err != nil {
			log.Critical("Error decoding PeerMap:", err)
		}
		log.Console(len(PeerMap.PL), "Peers loaded from persistant database")
		log.Trace("PeerMap loaded:", PeerMap.PL)
	} else {
		log.Debug("PeerMap database was empty")
	}

	// Do other initialization work
	go CronHourly()
	go CronDaily()

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
	err := FlushPeerMap()
	if err != nil {
		log.Error("Error in FlushPeerMap while shutting down:", err)
	}
	//  Yeah, yeah... the defer close might take effect before the caller exits, maybe...
	peersDB.Close()	
	c <- "Database shutdown complete"
}

func dbStatus() string {
	stats := peersDB.Stats()
	response := fmt.Sprintf("PeerInfo.bolt Reads: %d, Writes: %d, Nodes: %d, Rebalances: %d  Splits: %d",
		stats.TxN, stats.TxStats.Write, stats.TxStats.NodeCount, stats.TxStats.Rebalance, stats.TxStats.Split)
	return response
}

func FlushPeerMap() error {
	var encBuf bytes.Buffer

	err := peersDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(PeersBucket))
		if b == nil {
			log.Critical("Error accessing bucket:", PeersBucket)
		}
		log.Trace("FlushPeerMap:  Updating PeerMap")
		enc := gob.NewEncoder(&encBuf)
		err := enc.Encode(PeerMap.PL)
		if err != nil {
			log.Critical("FlushPeerMap: PeerMap Encode err:", err)
		}
		log.Trace("Gob len=", len(encBuf.String()))
		err = b.Put([]byte("PeerMap"), encBuf.Bytes())
		return err
	})
	if err != nil {
		log.Error("FlushPeerMap Update error:",err)
		return err
	}
	err = peersDB.Sync()
	if err != nil {
		log.Error("FlushPeerMap Sync error:",err)
		return err
	}
	log.Info("Information on", len(PeerMap.PL), "Peers flushed to persistant database")
	return nil
}
