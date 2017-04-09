// ArchIt logging routine
// Originally work created on 2/19/2017
//

package util

import (
	"bytes"
	"github.com/goarchit/archit/log"
	"net/http"
	"time"
)

func GetExtIP() string {
	// Go fetch our External IP address
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
