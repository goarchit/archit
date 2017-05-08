// Weird little function that is part of GoRPC, just having it defined gets it called
package farmer

import (
	"github.com/goarchit/archit/log"
	"github.com/valyala/gorpc"
	"errors"
	"io"
	"time"
)

var LastRemoteAddr string
var LastConnectTime time.Time

func newOnConnectFunc() gorpc.OnConnectFunc {
	return func(remoteAddr string, rwc io.ReadWriteCloser) (io.ReadWriteCloser, error) {
		now := time.Now()
		log.Console("NewOnConnect from", remoteAddr)
		if remoteAddr == LastRemoteAddr {
			if now.Sub(LastConnectTime) < ( 10 * time.Second) {
				return nil,errors.New("Spammer!")
			}
		}
		LastRemoteAddr = remoteAddr
		LastConnectTime = time.Now()
		return rwc, nil
	}
}
