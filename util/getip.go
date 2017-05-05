// ArchIt logging routine
// Originally work created on 2/19/2017
//

package util

import (
	"bytes"
	"github.com/goarchit/archit/log"
	"net"
	"net/http"
	"strings"
	"time"
)

func GetExtIP() string {
	// Go fetch our External IP address, WITHOUT PORT
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		log.Console("Problem getting PublicIP Address, retrying in 1 second")
		time.Sleep(1 * time.Second)
		resp, err = http.Get("http://myexternalip.com/raw")
		if err != nil {
			log.Console("Retrying again in 3 seconds")
			time.Sleep(3 * time.Second)
			resp, err = http.Get("http://myexternalip.com/raw")
			if err != nil {
				log.Console("Retrying one last time in 7 seconds")
				time.Sleep(7 * time.Second)
				resp, err = http.Get("http://myexternalip.com/raw")
				if err != nil {
					log.Console("Internet access is unreliable, shutting down")
					log.Critical(err)
				}
			}
		}

	}
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	s := buf.String()
	return s[0 : len(s)-1]
}

//function to get the public ip address
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Critical("Error getting outbound IP when connecting to 8.8.8.8", err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	return localAddr[0:idx]
}
