// dnsseed simply gathers the list of seed servers, stores them. and deteermines if we are one
// of them.

package util

import (
	"github.com/goarchit/archit/log"
	"github.com/valyala/gorpc"
	"net"
	"time"
	"math/rand"
)

var aliveDNSes []string

func Dnsseed() {
	var err error

	DNSSeeds, err = net.LookupHost("dnsseed.goarchit.online")
	aliveDNSes = make([]string,0)
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
			WG.Done()
		} else {
			go dnsalive(i,v)
		}
	}
	if SeedMode && !IAmASeed {
		log.Critical("Sorry, your public IP",PublicIP,"is not a registered seed node!")
	}
	WG.Wait()

	log.Debug(len(aliveDNSes),"DNSes found alive")
	rand.Seed(time.Now().UnixNano())
	MyDNSServerIP = aliveDNSes[rand.Intn(len(aliveDNSes)-1)]+SeedPortBase
	log.Console("Associating ourselves with DNSSeed",MyDNSServerIP)
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
		// Add mutexes just to be safe
		Mutex.Lock()
		// Add to current up list
		aliveDNSes = append(aliveDNSes, v)
		Mutex.Unlock()
		found = true
		log.Info("DNSSeed",v,"is alive.")
	}
	if !found && !IAmASeed {
		log.Critical("No DNSSeeds are apparently active.  Sorry.")
	}
}
