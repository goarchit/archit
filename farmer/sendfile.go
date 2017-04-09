package farmer

import (
	"bufio"
	"encoding/binary"
	"github.com/goarchit/archit/log"
	"io"
	"net"
	"os"
	"time"
)

func SendFile(address string, filename string) {
	starttime := time.Now().Unix()
	file, err := os.Open(filename)
	if err != nil {
		log.Critical(err)
	}
	defer file.Close()

	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Critical(err)
	}
	defer conn.Close()

	filenameSize := int64(len(filename))

	err = binary.Write(conn, binary.LittleEndian, filenameSize)
	if err != nil {
		log.Critical(err)
		return
	}

	_, err = io.WriteString(conn, filename)
	if err != nil {
		log.Error(err)
		return
	}

	stat, _ := file.Stat()
	err = binary.Write(conn, binary.LittleEndian, stat.Size())
	if err != nil {
		log.Error(err)
		return
	}

	br := bufio.NewReader(file)
	bw := bufio.NewWriter(conn)
	defer bw.Flush()

	_, err = io.CopyN(bw, br, stat.Size())
	if err != nil {
		log.Error(err)
		return
	}

	log.Console("File sent: ", filename, "to", address, "in", time.Now().Unix()-starttime, "seconds")
}
