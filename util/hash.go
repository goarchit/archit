// Hash.go, based the the goarchit/sample/hashfile.go routine
// Black2B is significantly faster than other hashes for our need
// Routine is called per 1GB slice and to create an authentication hash for each 32mb
// encoded slice
package util

import (
	"github.com/goarchit/archit/log"
        "github.com/minio/blake2b-simd"
	"io"
)

func HashString(source string) ([]byte, error) {
	var result []byte

	hash, err := blake2b.New(nil)
	if err != nil {
		log.Critical("Blake2B Hashing Error: ", err)
	}
	_, err = io.WriteString(hash, source)
	if err != nil {
		return result, err
	}

	return hash.Sum(result), nil
}
