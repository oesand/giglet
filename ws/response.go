package ws

import (
	"github.com/oesand/giglet"
	"github.com/oesand/giglet/internal/utils/stream"
	"github.com/oesand/giglet/specs"
	"net"
	"strings"
)

func UpgradeResponse(req giglet.Request, conf *Conf, handler Handler) giglet.Response {
	if req.Method() != specs.HttpMethodGet {
		return giglet.TextResponse("websocket: upgrading required request method - GET", specs.ContentTypePlain, func(response giglet.Response) {
			response.SetStatusCode(specs.StatusCodeMethodNotAllowed)
		})
	} else if !strings.EqualFold(req.Header().Get("Connection"), "Upgrade") {
		return giglet.TextResponse("websocket: 'Upgrade' token not found in 'Connection' header", specs.ContentTypePlain, func(response giglet.Response) {
			response.SetStatusCode(specs.StatusCodeBadRequest)
		})
	} else if !strings.EqualFold(req.Header().Get("Upgrade"), "websocket") {
		return giglet.TextResponse("websocket: 'websocket' token not found in 'Upgrade' header", specs.ContentTypePlain, func(response giglet.Response) {
			response.SetStatusCode(specs.StatusCodeBadRequest)
		})
	} else if req.Header().Get("Sec-Websocket-Version") != "13" {
		return giglet.TextResponse("websocket: supports only websocket 13 version", specs.ContentTypePlain, func(response giglet.Response) {
			response.SetStatusCode(specs.StatusCodeNotImplemented)
		})
	}

	challengeKey := req.Header().Get("Sec-Websocket-Key")
	if len(challengeKey) == 0 {
		return giglet.TextResponse("websocket: not a websocket handshake: `Sec-WebSocket-Key' header is missing or blank", specs.ContentTypePlain, func(response giglet.Response) {
			response.SetStatusCode(specs.StatusCodeBadRequest)
		})
	}
	if conf == nil {
		conf = &Conf{}
	}
	req.Hijack(func(conn net.Conn) {
		reader := stream.DefaultBufioReaderPool.Get(conn)
		defer stream.DefaultBufioReaderPool.Put(reader)

		wsConn := &wsConn{
			request: req,
			conn:    conn,
			reader:  *reader,
			conf:    *conf,
		}

		handler(wsConn)
		wsConn.dead = true
	})

	return giglet.EmptyResponse(specs.ContentTypeUndefined, func(resp giglet.Response) {
		resp.SetStatusCode(specs.StatusCodeSwitchingProtocols)
		resp.Header().Set("Upgrade", "websocket")
		resp.Header().Set("Connection", "Upgrade")
		resp.Header().Set("Sec-WebSocket-Accept", ComputeAcceptKey(challengeKey))
		if conf.EnableCompression {
			if ext := req.Header().Get("Sec-WebSocket-Extensions"); len(ext) > 0 && strings.Contains(ext, "permessage-deflate") {

				resp.Header().Set("Sec-WebSocket-Extensions", "permessage-deflate; server_no_context_takeover; client_no_context_takeover")
			}
		}
	})
}
