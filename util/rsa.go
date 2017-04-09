// RSA.go - A slightly modified version of:
// https://gist.github.com/sdorra/1c95de8cb80da31610d2ad767cd6f251
// Thanks go to the author of that page
/*
 * Genarate rsa keys.
 */

package util

import (
	"github.com/goarchit/archit/log"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/gob"
	"encoding/pem"
	"os"
)

func RSACheck(private, public string) error {
	if _, err := os.Stat(private+".key"); os.IsNotExist(err) {
		log.Debug(private+".key","file missing, recreating"
		cleanup(private, public)
		return err
	} 
	if _, err := os.Stat(private+".pem"); os.IsNotExist(err) {
		log.Debug(private+".pem","file missing, recreating"
		cleanup(private, public)
		return err
	} 
	if _, err := os.Stat(public+".key"); os.IsNotExist(err) {
		log.Debug(public+".key","file missing, recreating"
		cleanup(private, public)
		return err
	}
	if _, err := os.Stat(public+".pem"); os.IsNotExist(err) {
		log.Debug(public+".pem","file missing, recreating"
		cleanup(private, public)
		return err
	}
	return nil
}

func cleanup(private, public string) {
	// Remove all security files if they exist, ignore errors if they don't
	log.Debug("Some RSA files missing, clearing them all!")
	os.Remove(private+".key")
	os.Remove(private+".pem")
	os.Remove(public+".key")
	os.Remove(public+".pem")
}

func RSAGenerate(private, public string) {

	log.Debug("Generating RSA keys")
	reader := rand.Reader
	bitSize := 2048

	key, err := rsa.GenerateKey(reader, bitSize)
	checkError(err)

	publicKey := key.PublicKey

	saveGobKey(private+".key", key)
	savePEMKey(private+".pem", key)

	saveGobKey(public+".key", publicKey)
	savePublicPEMKey(public+".pem", publicKey)
}

func saveGobKey(fileName string, key interface{}) {
	outFile, err := os.Create(fileName)
	checkError(err)
	defer outFile.Close()

	encoder := gob.NewEncoder(outFile)
	err = encoder.Encode(key)
	checkError(err)
}

func savePEMKey(fileName string, key *rsa.PrivateKey) {
	outFile, err := os.Create(fileName)
	checkError(err)
	defer outFile.Close()

	var privateKey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(outFile, privateKey)
	checkError(err)
}

func savePublicPEMKey(fileName string, pubkey rsa.PublicKey) {
	asn1Bytes, err := asn1.Marshal(pubkey)
	checkError(err)

	var pemkey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	pemfile, err := os.Create(fileName)
	checkError(err)
	defer pemfile.Close()

	err = pem.Encode(pemfile, pemkey)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		log.Critical("Fatal error generating Security files", err.Error())
	}
}
