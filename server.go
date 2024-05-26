package giglet

import (
	"crypto/tls"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type Server struct {
	Handler RequestHandler
	Logger *log.Logger
	
	// ReadTimeout is the maximum duration for reading the entire
	// request, including the body. A zero or negative value means
	// there will be no timeout.
	ReadTimeout time.Duration

	// WriteTimeout is the maximum duration before timing out
	// writes of the response. It is reset whenever a new
	// request's header is read. Like ReadTimeout, it does not
	// let Handlers make decisions on a per-request basis.
	// A zero or negative value means there will be no timeout.
	WriteTimeout time.Duration

	// IdleTimeout is the maximum amount of time to wait for the
	// next request when keep-alives are enabled. If IdleTimeout
	// is zero, the value of ReadTimeout is used. If both are
	// zero, there is no timeout.
	IdleTimeout time.Duration

	TLSConfig *tls.Config

	nextProtos map[string]NextProtoHandler
	mutex sync.Mutex
	isShuttingdown atomic.Bool
	listenerTrack  sync.WaitGroup
}

func (server *Server) logger() *log.Logger {
	if server.Logger != nil {
		return server.Logger
	}
	return log.Default()
}

func (server *Server) handshakeTimeout() time.Duration {
	if server.ReadTimeout > 0 {
		return server.ReadTimeout
	} else if server.WriteTimeout > 0 {
		return server.WriteTimeout
	}
	return 0
}

func (s *Server) NextProto(key string, handler NextProtoHandler) {
	if s.nextProtos == nil {
		s.nextProtos = map[string]NextProtoHandler{}
	}
	if s.TLSConfig == nil {
		s.TLSConfig = &tls.Config{}
	}

	s.TLSConfig.NextProtos = append(s.TLSConfig.NextProtos, key)
	s.nextProtos[key] = handler
}

func (server *Server) ListenAndServe(addr string) error {
	lst, err := net.Listen("tcp4", addr)
	if err != nil {
		return err
	}
	return server.Serve(lst)
}
