// These functions are specific to those needed to send and receive files. 
// Originally work created on 3/19/2017
//
// Note the ../farmer/db.go routine handles all peerinfo data
//

package db

import (
	"bytes"
	"compress/flate"
	"encoding/gob"
	"fmt"
	"github.com/boltdb/bolt"
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

const FileBucket = "FileBucket"

var FileInfo DBRecord

var fileDB *bolt.DB

func Open() {

	var err error

	fileDB, err = bolt.Open(util.DBDir+"/FileInfo.bolt", 0600,
		&bolt.Options{Timeout: 2 * time.Second})
	if err != nil {
		log.Critical(err)
	}
	// Handle buckets
	err = fileDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(FileBucket))
		if err != nil {
			log.Critical("Error creating bucket:", FileBucket)
		}
		return nil
	})

	// Do other initialization work
	FileInfo.Slices = make(DBSlice)
}

func Close() {
	fileDB.Close()
}

func FileUpdate() {
	var encBuf, flateBuf bytes.Buffer
	err := fileDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(FileBucket))
		if b == nil {
			log.Critical("Error accessing bucket:", FileBucket)
		}
		key := fmt.Sprintf("%s@%s", FileInfo.Filename, FileInfo.UploadTime.Local())
		log.Trace("FileInfo.Bolt:  Updating key", key)
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
		log.Critical("FileInfo.bolt FileUpdate() error: ", err)
	}
}

func Status() string {
	stats := fileDB.Stats()
	response := fmt.Sprintf("FileInfo.bolt Reads: %d, Writes: %d, Nodes: %d, Rebalances: %d  Splits: %d\n",
		stats.TxN, stats.TxStats.Write, stats.TxStats.NodeCount, stats.TxStats.Rebalance, stats.TxStats.Split)
	return response
}
