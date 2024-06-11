package upgrading

import (
	"giglet"
	"giglet/safe"
	"giglet/specs"
	"net"
	"strings"
)

func UpgradeWebSocket(req giglet.Request, enableCompression bool, handler WebSocketHandler) giglet.Response {
	if req.Method() != specs.HttpMethodGet {
		return (&giglet.TextResponse{
			Text: "websocket: upgrading required request method - GET",
		}).SetStatusCode(specs.StatusCodeMethodNotAllowed)
	} else if !strings.EqualFold(req.Header().Get("Connection"), "Upgrade") {
		return (&giglet.TextResponse{
			Text: "websocket: 'Upgrade' token not found in 'Connection' header",
		}).SetStatusCode(specs.StatusCodeBadRequest)
	} else if !strings.EqualFold(req.Header().Get("Upgrade"), "websocket") {
		return (&giglet.TextResponse{
			Text: "websocket: 'websocket' token not found in 'Upgrade' header",
		}).SetStatusCode(specs.StatusCodeBadRequest)
	} else if req.Header().Get("Sec-Websocket-Version") != "13" {
		return (&giglet.TextResponse{
			Text: "websocket: supports only websocket 13 version",
		}).SetStatusCode(specs.StatusCodeNotImplemented)
	}
	
	challengeKey := req.Header().Get("Sec-Websocket-Key")
	if len(challengeKey) == 0 {
		return (&giglet.TextResponse{
			Text: "websocket: not a websocket handshake: `Sec-WebSocket-Key' header is missing or blank",
		}).SetStatusCode(specs.StatusCodeBadRequest)
	}
	req.Hijack(func(conn net.Conn) {
		handler(&WebSocketConn{
			request: req,
			conn: conn,
		})
	})
	return &WebSocketUpgradeResponse{
		ChallengeKey: challengeKey,
		EnableCompression: enableCompression,
	}
}

type WebSocketUpgradeResponse struct {
	_          safe.NoCopy

	ChallengeKey string
	EnableCompression bool
	header     *specs.Header
}

func (*WebSocketUpgradeResponse) StatusCode() specs.StatusCode {
	return specs.StatusCodeSwitchingProtocols
}

func (resp *WebSocketUpgradeResponse) Header() *specs.Header {
	if resp.header == nil {
		resp.header = &specs.Header{}

		resp.header.Set("Upgrade", "websocket")
		resp.header.Set("Connection", "Upgrade")
		resp.header.Set("Sec-WebSocket-Accept", specs.ComputeWebSocketAcceptKey(resp.ChallengeKey))
		if resp.EnableCompression {
			if ext :=resp.header.Get("Sec-WebSocket-Extensions"); 
				len(ext) > 0 && strings.Contains(ext, "permessage-deflate") {

				resp.header.Set("Sec-WebSocket-Extensions", "permessage-deflate; server_no_context_takeover; client_no_context_takeover");
			}
		}
	}
	return resp.header
}
