package upgrading

import "errors"

type WebSocketHandler func(conn *WebSocketConn)

var (
	ErrorInvalidWebsocketFrameType                   = errors.New("websocket: invalid frame type")
	ErrorWebsocketFrameSizeExceed                   = errors.New("websocket: frame size exceed")

	// Frame header byte 0 bits from Section 5.2 of RFC 6455
	webSocketFinalBit byte = 1 << 7
	webSocketRsv1Bit byte = 1 << 6
	webSocketRsv2Bit byte = 1 << 5
	webSocketRsv3Bit byte = 1 << 4

	// Frame header byte 1 bits from Section 5.2 of RFC 6455
	webSocketMaskBit byte = 1 << 7
)
