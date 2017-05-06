// dnsseed simply gathers the list of seed servers, stores them. and deteermines if we are one
// of them.

package util

import (
	"github.com/goarchit/archit/log"
	"github.com/valyala/gorpc"
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
			log.Console("We are a registered seed node!")
		}
		go dnsalive(i,v)
	}
	if SeedMode && !IAmASeed {
		log.Critical("Sorry, your public IP",PublicIP,"is not a registered seed node!")
	}
	WG.Wait()
}

func dnsalive(i int, v string) {
	var found bool
	defer WG.Done()
        //  Don't call yourself!
        pip,_,err := net.SplitHostPort(PublicIP)
        if err != nil {
                log.Critical("Error splitting PublicIP",err)
        }
        log.Trace("Checking IP",v,"against our PublicIP",pip,"?")
	// Don't call yourself
        if v == pip {
                log.Trace("Skipping talking to ourselves as unhealthy")
                return
        }
	serverIP := v+SeedPortBase
	c := gorpc.NewTCPClient(serverIP)
	c.Start()
	defer c.Stop()
        d := gorpc.NewDispatcher()
        d.AddFunc("Ping", func() {})
        dc := d.NewFuncClient(c)
        _, err = dc.CallTimeout("Ping", nil, time.Second*5)
	if err == nil {
	//	Mutex.Lock()
		// Any seed will do, so it doesn't matter who the last one is
		MyDNSServerIP = serverIP
	//	Mutex.Unlock()
		found = true
		log.Info("DNSSeed",v,"is alive.")
	}
	if !found && !IAmASeed {
		log.Critical("No DNSSeeds are apparently active.  Sorry.")
	}
}
