package engine

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

type sendPacketRequest struct {
	sid    string
	packet *packet
}

type getPacketRequest struct {
	sid string
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
