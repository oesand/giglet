package giglet

import (
	"crypto/tls"
	"errors"
	"log"
	"net"
	"slices"
	"sync"
	"sync/atomic"
	"time"
)

type Server struct {
	// Handler to invoke 
	Handler RequestHandler

	Logger *log.Logger

	Debug bool
	
	// Server name for sending in response headers.
	ServerName string

	// ReadTimeout is the maximum duration for reading the entire
	// request, including the body. A zero or negative value means
	// there will be no timeout.
	ReadTimeout time.Duration

	// WriteTimeout is the maximum duration before timing out
	// writes of the response. A zero or negative value means
	// there will be no timeout.
	WriteTimeout time.Duration

	// TLSConfig optionally provides a TLS configuration 
	TLSConfig *tls.Config

	// ContentMaxSizeBytes controls the maximum number of bytes the
	// server will read parsing the request header's keys and
	// values, including the request line and the request body.
	// If zero, DefaultContentMaxSizeBytes is used.
	ContentMaxSizeBytes int64

	nextProtos map[string]NextProtoHandler
	isShuttingdown atomic.Bool
	listenerTrack  sync.WaitGroup
}

func (server *Server) logger() *log.Logger {
	if server.Logger != nil {
		return server.Logger
	}
	return log.Default()
}

func (server *Server) applyReadTimeout(conn net.Conn) {
	if server.ReadTimeout > 0 {
		conn.SetReadDeadline(time.Now().Add(server.ReadTimeout))
	}
}

func (server *Server) applyWriteTimeout(conn net.Conn) {
	if server.WriteTimeout > 0 {
		conn.SetWriteDeadline(time.Now().Add(server.WriteTimeout))
	}
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
	if server.isShuttingdown.Load() {
		return ErrorServerClosed
	} else if addr == "" {
		addr = ":http"
	}
	lst, err := net.Listen("tcp4", addr)
	if err != nil {
		return err
	}
	return server.Serve(lst)
}

func (server *Server) ListenAndServeTLS(addr, certFile, keyFile string) error {
	if server.isShuttingdown.Load() {
		return ErrorServerClosed
	} else if addr == "" {
		addr = ":http"
	}
	lst, err := net.Listen("tcp4", addr)
	if err != nil {
		return err
	}
	return server.ServeTLS(lst, certFile, keyFile)
}

func (server *Server) ListenAndServeTLSRaw(addr string, cert *tls.Certificate) error {
	if server.isShuttingdown.Load() {
		return ErrorServerClosed
	} else if addr == "" {
		addr = ":http"
	}
	lst, err := net.Listen("tcp4", addr)
	if err != nil {
		return err
	}
	return server.ServeTLSRaw(lst, cert)
}

func (srv *Server) ServeTLS(lst net.Listener, certFile, keyFile string) error {
	if len(certFile) == 0 || len(keyFile) == 0 {
		return errors.New("unknown certificate source")
	}
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}
	return srv.ServeTLSRaw(lst, &cert)
}

func (srv *Server) ServeTLSRaw(lst net.Listener, cert *tls.Certificate) error {
	// Setup HTTP/2 before srv.Serve, to initialize srv.TLSConfig
	// before we clone it and create the TLS Listener.
	// if err := srv.setupHTTP2_ServeTLS(); err != nil {
	// 	return err
	// }

	var config *tls.Config
	if srv.TLSConfig != nil {
		config = srv.TLSConfig.Clone()
	} else {
		config = &tls.Config{}
	}

	if !slices.Contains(config.NextProtos, "http/1.1") {
		config.NextProtos = append(config.NextProtos, "http/1.1")
	}

	configHasCert := len(config.Certificates) > 0 || config.GetCertificate != nil
	if !configHasCert || cert != nil {
		config.Certificates = make([]tls.Certificate, 1)
		config.Certificates[0] = *cert
	}

	listener := tls.NewListener(lst, config)
	return srv.Serve(listener)
}

func (srv *Server) Shutdown() {
	srv.isShuttingdown.Store(true)
	srv.listenerTrack.Wait()
}
