package specs

import (
	"crypto/sha1"
	"encoding/base64"
)

type WebSocketFrame uint16

const (
	WebSocketTextFrame   WebSocketFrame = 1
	WebSocketBinaryFrame WebSocketFrame = 2
	WebSocketCloseFrame  WebSocketFrame = 8
	WebSocketPingFrame   WebSocketFrame = 9
	WebSocketPongFrame   WebSocketFrame = 10
)

var (
	WebSocketCloseNormal              = &StatusCode{Code: 1000, Detail: []byte("(normal)")}
	WebSocketCloseGoingAway           = &StatusCode{Code: 1001, Detail: []byte("(going away)")}
	WebSocketCloseProtocolError       = &StatusCode{Code: 1002, Detail: []byte("(protocol error)")}
	WebSocketCloseUnsupportedData     = &StatusCode{Code: 1003, Detail: []byte("(unsupported data)")}
	WebSocketCloseNoStatusReceived    = &StatusCode{Code: 1005, Detail: []byte("(no status)")}
	WebSocketCloseAbnormal            = &StatusCode{Code: 1006, Detail: []byte("(abnormal closure)")}
	WebSocketCloseInvalidPayloadData  = &StatusCode{Code: 1007, Detail: []byte("(invalid payload data)")}
	WebSocketClosePolicyViolation     = &StatusCode{Code: 1008, Detail: []byte("(policy violation)")}
	WebSocketCloseMessageTooBig       = &StatusCode{Code: 1009, Detail: []byte("(message too big)")}
	WebSocketCloseMandatoryExtension  = &StatusCode{Code: 1010, Detail: []byte("(mandatory extension missing)")}
	WebSocketCloseInternalServerError = &StatusCode{Code: 1011, Detail: []byte("(internal server error)")}
	WebSocketCloseServiceRestart      = &StatusCode{Code: 1012, Detail: []byte("(service restart)")}
	WebSocketCloseTryAgainLater       = &StatusCode{Code: 1013, Detail: []byte("(try again later)")}
	WebSocketCloseTLSHandshake        = &StatusCode{Code: 1015, Detail: []byte("(tls handshake error)")}
)

func ComputeWebSocketAcceptKey(challengeKey []byte) string {
	h := sha1.New() // (CWE-326) -- https://datatracker.ietf.org/doc/html/rfc6455#page-54
	h.Write(challengeKey)
	h.Write(websocketAcceptBaseKey)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
