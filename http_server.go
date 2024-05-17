package giglet

import (
	"bufio"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type Server struct {
	Handler RequestHandler

	Logger *log.Logger
	
	TCPKeepAlivePeriod time.Duration

	ReadBufferSize int

	mutex sync.Mutex
}


func (server *Server) logger() *log.Logger {
	if server.Logger != nil {
		return server.Logger
	}
	return log.Default()
}
func (server *Server) accept(lst net.Listener) (net.Conn, error) {
	for {
		conn, err := lst.Accept()

		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				server.logger().Printf("timeout error when accepting new connections: %v", netErr)
				continue
			}
			return nil, err
		} else if tcp, ok := conn.(*net.TCPConn); ok && server.TCPKeepAlivePeriod > 0 {
			err := tcp.SetKeepAlive(true);
			if err != nil {
				err = tcp.SetKeepAlivePeriod(server.TCPKeepAlivePeriod)
			}
			if err != nil {
				tcp.Close()
				return nil, err
			}
		}

		return conn, nil
	}
}

func (server *Server) ListenAndServe(addr string) error {
	lst, err := net.Listen("tcp4", addr)
	if err != nil {
		return err
	}
	return server.Serve(lst)
}

func (server *Server) Serve(listener net.Listener) error {
	for {
		conn, err := server.accept(listener)

		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		go server.work(conn)
	}
}

func (server *Server) work(conn net.Conn) {
	reader := bufio.NewReaderSize(conn, server.ReadBufferSize)
	
	// request := HttpRequest{
	// 	conn: conn,
	// }
	
	for {

		reader.ReadLine()

		buffer, _ := reader.Peek(1)
		if len(buffer) == 0 {
			conn.Close()
		}

		

		
	}
}
