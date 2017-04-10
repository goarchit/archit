package farmer

import (
	"github.com/goarchit/archit/log"
	"github.com/valyala/gorpc"
	"io"
	"sync"
)

var Connect sync.Mutex	
var RemoteAddr string

func newOnConnectFunc() gorpc.OnConnectFunc {
	return func(remoteAddr string, rwc io.ReadWriteCloser) (io.ReadWriteCloser, error) {
		log.Console("Connection from", remoteAddr)
		Connect.Lock()	// Wanring, must be unlocked by called routine
		RemoteAddr = remoteAddr
		return rwc, nil
	}
}
