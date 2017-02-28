package archcrypt_test
import (
	"github.com/goarchit/archit/archcrypt"
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/log"
	"fmt"
	"testing"
)

func TestEncrypt(t *testing.T) {
	log.Setup(4,"test.log",2, true)
	config.DerivedKey = fmt.Sprintf("%32s","Stupid test key")
	log.Debug("DerivedKey = ", config.DerivedKey)
	archcrypt.DoOnce()
	log.Debug("DoOnce() called")
	s := fmt.Sprintf("%48s","The fox jumped over some rocks, tripped, and died.")
	log.Debug("s =",s)
	cipherblock := archcrypt.Encrypt([]byte(s))
	log.Debug("cipherblock = ",cipherblock)
	plaintext := archcrypt.Decrypt(cipherblock)
	log.Debug("plaintext = ",plaintext)
}

