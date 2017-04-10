package farmer

import (
	"github.com/goarchit/archit/log"
	"github.com/valyala/gorpc"
	"io"
)

var RemoteAddr string

func newOnConnectFunc() gorpc.OnConnectFunc {
	return func(remoteAddr string, rwc io.ReadWriteCloser) (io.ReadWriteCloser, error) {
		log.Console("Connection from", remoteAddr)
		RemoteAddr = remoteAddr
		return rwc, nil
	}
}
