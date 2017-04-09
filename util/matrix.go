// Matrix - build the decoding matricing
// Originally work created on 3/14/2017
//

package util

import (
	"github.com/goarchit/archit/log"
)

var Matrix [32 + MaxRaptor][32]int
var MatrixCount [32 + MaxRaptor]int

func BuildMatrix() {
	log.Trace("Matrix:  Raptor=", Raptor, "MaxRaptor=", MaxRaptor)
	for i := 0; i < 32+Raptor; i++ {
		// Main data blocks
		if i < 32 {
			for j := i; j < i+Raptor; j++ {
				if j < 32 {
					Matrix[i][j] = 1
					MatrixCount[i]++
					log.Trace("Matrix[", i, "][", j, "]=1")
				} else {
					Matrix[i][j-32] = 1
					MatrixCount[i]++
					log.Trace("Matrix[", i, "][", j-32, "]=1")
				}
			}
			// Redundancy blocks
		} else {
			Matrix[i][i-32] = 1
			MatrixCount[i]++
			log.Trace("Matrix[", i, "][", i-32, "]=1")
		}
		log.Trace("MatrixCount[", i, "]=", MatrixCount[i])
	}
}
