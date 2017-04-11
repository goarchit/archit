// dnsseed simply gathers the list of seed servers, stores them. and deteermines if we are one
// of them.

package util

import (
	"github.com/goarchit/archit/log"
	"net"
	"time"
)

func Dnsseed() {
	var err error

	DNSSeeds, err = net.LookupHost("dnsseed.goarchit.online")
	if err != nil {
                log.Critical("Failure to lookup dnsseed.goarchit.online")
        }
        log.Debug("DNSseed resolved to",len(DNSSeeds),"seeds:", DNSSeeds)
	if len(PublicIP) < 9 {   // len(1.2.3.4:1) == 9
		log.Critical("Dnsseed - PublicIP is obviously wrong:",PublicIP)
	}
	WG.Add(len(DNSSeeds))
	for i, v := range DNSSeeds {
		ip := v+SeedPortBase
		if ip == PublicIP {
			IAmASeed = true
			log.Console("We are a seed node!  No storage farming allowed.")
		}
		go dnsalive(i,v)
	}
	WG.Wait()
}

func dnsalive(i int, v string) {
	var found bool
	defer WG.Done()
	serverIP := v+SeedPortBase
	con, err := net.DialTimeout("tcp", serverIP, time.Second*10)
	if err == nil {
		Mutex.Lock()
		// Any seed will do, so it doesn't matter who the last one is
		MyDNSServerIP = serverIP
		Mutex.Unlock()
		con.Close()
		found = true
		log.Info("DNSSeed",v,"is alive.")
	}
	if !found && !IAmASeed {
		log.Critical("No DNSSeeds are apparently active.  Sorry.")
	}
}
