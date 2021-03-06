package farmer

import (
	"github.com/goarchit/archit/log"
	"encoding/json"
	"strconv"
)

func Status() string {
	walletCmd <- "status"
	response := <-walletCmd
	response += "\n"+strconv.Itoa(len(PeerMap.PL))+" Peers: [ "
	count := 0
	for _, v := range PeerMap.PL {
		count++
		// Limit screen output
		if count > 100 { 
			response += "..."
			break 
		}
		response += v.IPAddr + " "
	}
	response += "]\n\nInternal RPC Stats: "
	s, err := json.Marshal(intCmd.Stats)
	if err != nil {
		log.Critical("Farmer status() error Marshaling intCmd.Stats")
	}
	response += string(s)+"\n"
	response += "\nExternal RPC Stats: "
	s, err = json.Marshal(extCmd.Stats)
	if err != nil {
		log.Critical("Farmer status() error Marshaling extCmd.Stats")
	}
	response += string(s)+"\n"
	response += "]\n\nDatabase statuses since last open:\n"
	dbCmd <- "status"
	response += <-dbCmd
	log.Info("Farmer Status returning:\n", response)
	return response
}
