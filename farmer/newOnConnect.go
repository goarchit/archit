//  Try and prevent spamming by limiting the number of request to once every 10 seconds
//  (unless someeone else connects)
//  This function is only set for External connectionings in run.go
package farmer

import (
	"github.com/goarchit/archit/log"
	"github.com/valyala/gorpc"
	"io"
	"net"
	"time"
)

var LastRemoteAddr string
var LastConnectTime time.Time

func newOnConnectFunc() gorpc.OnConnectFunc {
	return func(remoteAddr string, rwc io.ReadWriteCloser) (io.ReadWriteCloser, error) {
		now := time.Now()
		log.Trace("NewOnConnect from", remoteAddr)
		ra, _, err := net.SplitHostPort(remoteAddr)
		if err != nil {
			log.Critical("newOnConnectFunc:  Error Splitting IP address",remoteAddr,":",err)
		}
		if ra == LastRemoteAddr {
			if now.Sub(LastConnectTime) < ( 10 * time.Second) {
				log.Trace("Anti-Spammer delay being added")
				time.Sleep(1958 * time.Millisecond)
			}
		} else {
			LastRemoteAddr = ra
		}
		LastConnectTime = time.Now()
		return rwc, nil
	}
}
