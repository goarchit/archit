// dnsseed simply gathers the list of seed servers, stores them. and deteermines if we are one
// of them.

package util

import (
	"github.com/goarchit/archit/log"
	"net"
)

func Dnsseed() {
	var err error
	
	DNSSeeds, err = net.LookupHost("dnsseed.goarchit.online")
	if err != nil {
                log.Critical("Failure to lookup dnsseed.goarchit.online")
        }
        log.Debug("DNSseed resolved to", DNSSeeds)
	if PublicIP == "" {
		PublicIP = GetExtIP()
		log.Debug("Dnsseed - didn't expect to have to get our public IP")
	}
	for _, v := range DNSSeeds {
		if v == PublicIP {
			IAmASeed = true
			log.Console("We are a seed node!  Behave!!!")
		}	
	}

}
