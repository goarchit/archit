package farmer

import (
	"github.com/goarchit/archit/log"
)

func Status () string {
	walletCmd <- "status"
	response := <- walletCmd
	response += "\nDatabase Statuses:\n"
	dbCmd <- "status"
	response += <- dbCmd
	response += "\nInternal RPC server status:\n"
//	stats := intCmd.ConnStats.Snapshot()
//	str, _ := json.Marshal(stats)
//	log.Console(string(str))
	response += "\nExternal RPC server status:\n"
//	stats = extCmd.ConnStats.Snapshot()
//	str, _ = json.Marshal(stats)
//	log.Console(string(str))
	log.Info("Farmer Status returning:\n",response)
	return response
}
