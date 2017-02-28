// ArchIt encrypting and decrypting routine
// Originally work created on 1/25/2017
//

package archcrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/log"
)

/* NOT FINISHED CODING */
var KeyBlock cipher.Block

func DoOnce() {
	KeyBlock, err := aes.NewCipher([]byte(config.DerivedKey))
	if err != nil {
		panic(err)
	}
	log.Debug("KeyBlock:",KeyBlock)
}

func Encrypt(plaintext []byte) []byte {

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Critical(err)
	}

	stream := cipher.NewCTR(KeyBlock, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.
	return ciphertext
}

func Decrypt(cipherblock []byte) []byte {

	// CTR mode is the same for both encryption and decryption, so we can
	// also decrypt that ciphertext with NewCTR.

	plaintext2 := make([]byte, len(cipherblock)-aes.BlockSize)
	iv := cipherblock[:aes.BlockSize]
	stream := cipher.NewCTR(KeyBlock, iv)
	stream.XORKeyStream(plaintext2, cipherblock[aes.BlockSize:])

	return plaintext2
}
