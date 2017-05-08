//  Filename manipulation routines initially written on 5/6/17
package farmer

import (

)

func FindNodeIP(nodeIP string) string {
	for k,v := range PeerMap.PL {
		if v.IPAddr == nodeIP {
			return k
		} 	
	}
	return ""
}
