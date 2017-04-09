// Include... common definitions and external variables
// Originally work created on 3/19/2017
//

package db

import (
	"bytes"
	"compress/flate"
	"encoding/gob"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/farmer"
	"github.com/goarchit/archit/log"
	"github.com/goarchit/archit/util"
	"io"
	"sync"
	"time"
)

type DBSlice map[string]struct {
	Raptor int
	Shard  [32 + util.MaxRaptor]struct {
		Wallet string // Wallet Address of the farmer
		Hmac   []byte
		Random [32]struct { // 32 random samples from a Shard
			Offset int
			Value  [4]byte
		}
	}
}

type DBRecord struct {
	Filename   []byte
	UploadTime time.Time
	Mutex      sync.Mutex
	Slices     DBSlice
}

const FileBucket string = "FileInfo"
const RepBucket string = "Reputation"

var FileInfo DBRecord

var db *bolt.DB

func Open() {
	var peerbytes bytes.Buffer
	var err error
	db, err = bolt.Open(config.Archit.DBDir+"/Archit.bolt", 0600,
		&bolt.Options{Timeout: 2 * time.Second})
	if err != nil {
		log.Critical(err)
	}
	// Handle buckets
	err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(FileBucket))
		if err != nil {
			log.Critical("Error creating bucket:", FileBucket)
		}
		b, err := tx.CreateBucketIfNotExists([]byte(RepBucket))
		if err != nil {
			log.Critical("Error creating bucket:", RepBucket)
		}
		// Load up the reputation matrix
		peerbytes.Write(b.Get([]byte("PeerMap")))
		return nil
	})
	if peerbytes.Len() != 0 {
		dec := gob.NewDecoder(&peerbytes)
		err := dec.Decode(&farmer.PeerMap.PL)
		if err != nil {
			log.Critical("Error decoding PeerMap:", err)
		}
		log.Debug("PeerMap successfully assigned")
		log.Trace("PeerMap loaded:",farmer.PeerMap.PL)
	}

	// Do other initialization work
	go CronHourly()
	go CronDaily()
	FileInfo.Slices = make(DBSlice)
}

func Close() {
	FlushPeerMap()
	db.Close()
}

func FileUpdate() {
	var encBuf, flateBuf bytes.Buffer
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(FileBucket))
		if b == nil {
			log.Critical("Error accessing bucket:", FileBucket)
		}
		key := fmt.Sprintf("%s@%s", FileInfo.Filename, FileInfo.UploadTime.Local())
		log.Trace("Archit.Bolt:  Updating key", key)
		enc := gob.NewEncoder(&encBuf)
		err := enc.Encode(FileInfo.Slices)
		if err != nil {
			log.Critical("DB: Encode err:", err)
		}
		log.Trace("Gob len=", len(encBuf.String()))
		fw, err := flate.NewWriter(nil, flate.BestCompression)
		if err != nil {
			log.Critical("DB: Flate err:", err)
		}
		fw.Reset(&flateBuf)
		_, err = io.WriteString(fw, encBuf.String())
		if err != nil {
			log.Critical("DB: WriteString err:", err)
		}
		fw.Close()
		log.Trace("Flate, BestCompression=", len(flateBuf.String()))

		err = b.Put([]byte(key), flateBuf.Bytes())
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Critical("Archit.bolt FileUpdate error: ", err)
	}
}

func Status() string {
	stats := db.Stats()
	response := fmt.Sprintf("Archit.bolt Reads: %d, Writes: %d, Nodes: %d, Rebalances: %d  Splits: %d\n",
		stats.TxN, stats.TxStats.Write, stats.TxStats.NodeCount, stats.TxStats.Rebalance, stats.TxStats.Split)
	return response
}

func FlushPeerMap() {
	var encBuf bytes.Buffer

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(RepBucket))
		if b == nil {
			log.Critical("Error accessing bucket:", RepBucket)
		}
		log.Trace("Archit.Bolt:  Updating PeerMap")
		enc := gob.NewEncoder(&encBuf)
		err := enc.Encode(farmer.PeerMap.PL)
		if err != nil {
			log.Critical("DB: PeerMap Encode err:", err)
		}
		log.Trace("Gob len=", len(encBuf.String()))
		err = b.Put([]byte("PeerMap"), encBuf.Bytes())
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Critical("Archit.bolt Put error updating PeerMap: ", err)
	}
	err = db.Sync()
	if err != nil {
		log.Critical("Archit.bolt error syncing database to disk: ", err)
	}
}
