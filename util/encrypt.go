// Encrypting - AES encryption of encoded blocks
// IV is stored at the END of each block, not the beginning
// Originally work created on 3/19/2017
//

package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/goarchit/archit/log"
	"io"
)

func Encrypt(tblock []byte, iv *[32 + MaxRaptor][aes.BlockSize]byte, i int) {
	defer WG.Done()
	log.Trace("Starting Encrypting block", i)
	block, err := aes.NewCipher(DerivedKey)
	if err != nil {
		log.Critical("Encrypt error: ", err)
	}
	if len(DerivedKey) < aes.BlockSize {
		log.Critical("Encyrpt error: AES Blocksize", aes.BlockSize, "greater then len(DerivedKey)", len(DerivedKey))
	}

	tiv := iv[i][:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, tiv); err != nil {
		log.Critical("Encrypt error: IV Generation:", err)
	}
	if len(tiv) != aes.BlockSize {
		log.Critical("Encrypt logic error, len(tiv) =", len(tiv))
	}

	stream := cipher.NewCTR(block, tiv)
	stream.XORKeyStream(tblock, tblock)
}
