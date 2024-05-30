package specs

import (
	"crypto/sha1"
	"encoding/base64"
	"giglet/safe"
	"strconv"
)

type WebSocketFrame uint16

const (
	WebSocketTextFrame   WebSocketFrame = 1
	WebSocketBinaryFrame WebSocketFrame = 2
	WebSocketCloseFrame  WebSocketFrame = 8
	WebSocketPingFrame   WebSocketFrame = 9
	WebSocketPongFrame   WebSocketFrame = 10
)

type WebSocketClose struct {
	_ safe.NoCopy
	
	Code uint16
	Detail []byte
}

func (status *WebSocketClose) AppendBytes(buffer []byte) []byte {
	buffer = strconv.AppendUint(buffer, uint64(status.Code), 10)
	buffer = append(buffer, ' ')
	buffer = append(buffer, status.Detail...)
	return buffer
}

var (
	WebSocketCloseNormal              = &WebSocketClose{Code: 1000, Detail: []byte("(normal)")}
	WebSocketCloseGoingAway           = &WebSocketClose{Code: 1001, Detail: []byte("(going away)")}
	WebSocketCloseProtocolError       = &WebSocketClose{Code: 1002, Detail: []byte("(protocol error)")}
	WebSocketCloseUnsupportedData     = &WebSocketClose{Code: 1003, Detail: []byte("(unsupported data)")}
	WebSocketCloseNoStatusReceived    = &WebSocketClose{Code: 1005, Detail: []byte("(no status)")}
	WebSocketCloseAbnormal            = &WebSocketClose{Code: 1006, Detail: []byte("(abnormal closure)")}
	WebSocketCloseInvalidPayloadData  = &WebSocketClose{Code: 1007, Detail: []byte("(invalid payload data)")}
	WebSocketClosePolicyViolation     = &WebSocketClose{Code: 1008, Detail: []byte("(policy violation)")}
	WebSocketCloseMessageTooBig       = &WebSocketClose{Code: 1009, Detail: []byte("(message too big)")}
	WebSocketCloseMandatoryExtension  = &WebSocketClose{Code: 1010, Detail: []byte("(mandatory extension missing)")}
	WebSocketCloseInternalServerError = &WebSocketClose{Code: 1011, Detail: []byte("(internal server error)")}
	WebSocketCloseServiceRestart      = &WebSocketClose{Code: 1012, Detail: []byte("(service restart)")}
	WebSocketCloseTryAgainLater       = &WebSocketClose{Code: 1013, Detail: []byte("(try again later)")}
	WebSocketCloseTLSHandshake        = &WebSocketClose{Code: 1015, Detail: []byte("(tls handshake error)")}
)

func ComputeWebSocketAcceptKey(challengeKey string) string {
	h := sha1.New() // (CWE-326) -- https://datatracker.ietf.org/doc/html/rfc6455#page-54
	h.Write([]byte(challengeKey))
	h.Write(websocketAcceptBaseKey)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
