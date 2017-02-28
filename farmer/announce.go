package farmer

import (
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/log"
	"github.com/goarchit/archit/util"
	"encoding/json"
	"net"
	"strconv"
	"strings"
)

func announce() {
	
	addresses, err := net.LookupHost("dnsseed.goarchit.online")
	if err != nil {
		log.Critical("Failure to lookup dnsseed.goarchit.online")
	}	
	log.Debug("DNSseed resolved to", addresses)
	for i := 0; i < len(addresses); i++ {
		seed := strings.SplitAfter(addresses[i], " ")
		log.Info("Found seed ",seed)
	}

	publicIP := util.GetExtIP()

	iAm := new(PeerInfo)
	iAm.Address = config.Archit.WalletAddr
	iAm.Detail.IpAddr = publicIP+":"+strconv.Itoa(config.Archit.PortBase)
	iAm.Detail.MacAddr = "Invalid"
	rifs := util.RoutedInterface("ip", net.FlagUp|net.FlagBroadcast)
        if rifs != nil {
		iAm.Detail.MacAddr = rifs.HardwareAddr.String()
        }
	s, _ := json.Marshal(iAm)
	log.Console("Farmer node startup complete!")
	log.Debug("whoAmI:",string(s))
}
