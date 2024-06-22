package engine

type Transport string

const (
	Polling      Transport = "polling"
	WebSocket    Transport = "websocket"
	WebTransport Transport = "webtransport"
)
