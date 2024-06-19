package upgrading

import (
	"errors"
	"giglet/safe"
)

type WebSocketHandler func(conn *WebSocketConn)

var bufioReaderPool = safe.BufioReaderPool{
	MaxSize: 128,
}

var (
	ErrorWebsocketInvalidFrameType  = errors.New("websocket: invalid frame type")
	ErrorWebsocketFrameSizeExceed   = errors.New("websocket: frame size exceed")
	ErrorWebsocketClosed            = errors.New("websocket: closed")

	ErrorWebsocketNoRsV1            = errors.New("websocket: rsv1 not implemented")
	ErrorWebsocketNoRsV2            = errors.New("websocket: rsv2 not implemented")
	ErrorWebsocketNoRsV3            = errors.New("websocket: rsv3 not implemented")
)

const (
	// Frame header byte 0 bits from Section 5.2 of RFC 6455
	websocketFinalBit byte = 1 << 7
	websocketRsv1Bit byte = 1 << 6
	websocketRsv2Bit byte = 1 << 5
	websocketRsv3Bit byte = 1 << 4

	// Frame header byte 1 bits from Section 5.2 of RFC 6455
	websocketMaskBit byte = 1 << 7

	websocketMaxFrameHeaderSize         = 2 + 8 + 4 // Fixed header + length + mask
	websocketMaxServiceFramePayloadSize = 125

	// minCompressionLevel     = -2 // flate.HuffmanOnly not defined in Go < 1.6
	// maxCompressionLevel     = flate.BestCompression
	// defaultCompressionLevel = 1
)

