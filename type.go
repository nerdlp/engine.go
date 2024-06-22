package engine

import "context"

type Transport string

const (
	Polling      Transport = "polling"
	WebSocket    Transport = "websocket"
	WebTransport Transport = "webtransport"
)

type QueryKey = string

const (
	qk_EIO       QueryKey = "EIO"
	qk_Transport QueryKey = "transport"
	qk_sid                = "sid"
)

type sendMessageRequest struct {
	sid         string
	messageType int
	message     []byte
}

type transportClient interface {
	Send(ctx context.Context, request *sendMessageRequest) error
}

type packetType = byte

const (
	open    packetType = 0
	close   packetType = 1
	ping    packetType = 2
	pong    packetType = 3
	message packetType = 4
	upgrade packetType = 5
	noop    packetType = 6
)
