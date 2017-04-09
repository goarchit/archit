// Encoding - xor blocks together for later decoding
// Originally work created on 3/14/2017
//

package util

import (
	"github.com/goarchit/archit/log"
)

func Encode(tblock *[(32 + MaxRaptor) * ShardLen]byte, block *[32 * ShardLen]byte, i int) {
	defer WG.Done()
	sliceOfTblock := tblock[i*ShardLen : (i+1)*ShardLen-1]
	sliceOfBlock := make([]byte, ShardLen, ShardLen)
	for j := 1; j < Raptor; j++ {
		off := (i + j) * ShardLen
		if i+j < 32 {
			log.Trace("XORing block", i, "with block", i+j)
			sliceOfBlock = block[off : off+ShardLen-1]
		} else {
			log.Trace("XORing block", i, "with tblock", i+j)
			sliceOfBlock = tblock[off : off+ShardLen-1]
		}
		XorBytes(sliceOfTblock, sliceOfBlock)
	}
}
