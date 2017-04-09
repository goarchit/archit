// hmac.go is located in db instead of util to avoid import cycle issues
package db

import (
	"github.com/goarchit/archit/log"
	"github.com/goarchit/archit/util"
	"math/rand"
	"time"
)

func HMAC(source string, i int) {
	var value [4]byte
	defer util.WG.Done()
	hash, err := util.HashString(source)
	if err != nil {
		log.Critical("HashString error: ", err)
	}
	log.Trace("Slice:", i, "Authentication hash =", hash)
	rand.Seed(time.Now().UnixNano())
	FileInfo.Mutex.Lock() // Protect Map activity
	defer FileInfo.Mutex.Unlock()
	s := FileInfo.Slices[util.SliceName]
	s.Shard[i].Hmac = hash
	for j := 0; j < 32; j++ {
		off := rand.Intn(util.ShardLen - 5) // 4 should do, but be safe
		s.Shard[i].Random[j].Offset = off
		for k := 0; k < 4; k++ {
			value[k] = source[off+k]
		}
		s.Shard[i].Random[j].Value = value
	}
	FileInfo.Slices[util.SliceName] = s
}
