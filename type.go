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
)
