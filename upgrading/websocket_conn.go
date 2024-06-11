package upgrading

import (
	"giglet"
	"giglet/specs"
	"io"
	"net"
)

type WebSocketConn struct {
	request giglet.Request
	conn net.Conn
	dead bool
}

func (conn *WebSocketConn) GetExtra(key string) any {
	return conn.request.GetExtra(key)
}

func (conn *WebSocketConn) SetExtra(key string, value any) {
	conn.request.SetExtra(key, value)
}

func (conn *WebSocketConn) RemoteAddr() net.Addr {
	return conn.request.RemoteAddr()
}

func (conn *WebSocketConn) Url() *specs.Url {
	return conn.request.Url()
}

func (conn *WebSocketConn) Header() *specs.ReadOnlyHeader {
	return conn.request.Header()
}

func (conn *WebSocketConn) Alive() bool {
	return !conn.dead
}

func (conn *WebSocketConn) WriteServiceFrame(frameType specs.WebSocketFrame, payload []byte) error {
	if !frameType.IsService() {
		return ErrorInvalidWebsocketFrameType
	} else if len(payload) > 125 {
		return ErrorWebsocketFrameSizeExceed
	}

	buf := make([]byte, 0, 140) // max header & frame size 
	buf = append(buf, byte(frameType) | webSocketFinalBit, byte(len(payload)))
	buf = append(buf, payload...)

	_, err := conn.conn.Write(buf)
	if err == io.EOF {
		conn.dead = true
	} else if frameType == specs.WebSocketCloseFrame {
		conn.dead = true
		err = io.EOF
		conn.conn.Close()
	}
	return err
}

func (conn *WebSocketConn) WriteCloseFrame(reason specs.WebSocketClose) error {
	return conn.WriteServiceFrame(specs.WebSocketCloseFrame, reason.Bytes())
}

