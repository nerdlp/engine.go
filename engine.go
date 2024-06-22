package engine

import (
	"net/http"
	"time"
)

type Option func(engine *Engine)

const (
	defaultPingInterval = 20000 * time.Millisecond
	defaultPingTimeout  = 25000 * time.Millisecond
	defaultMaxPayload   = 1 << 10 // 1 KB
)

var (
	defaultTransport = []Transport{Polling, WebSocket}
)

type Engine struct {
	// how many ms before sending a new ping packet
	pingInterval time.Duration
	// how many ms without a pong packet to consider the connection closed
	pingTimeout time.Duration
	// how many bytes or characters a message can be, before closing the session (to avoid DoS)
	maxPayload uint32
	// The low level transports enabled
	transports []Transport
	// allows to upgrade transport
	allowUpgrades bool
	// path
	path string
}

func New(opts ...Option) *Engine {
	engine := &Engine{
		pingInterval:  defaultPingInterval,
		pingTimeout:   defaultPingTimeout,
		maxPayload:    defaultMaxPayload,
		transports:    defaultTransport,
		allowUpgrades: true,
	}
	for _, opt := range opts {
		opt(engine)
	}

	return engine
}

func (e *Engine) Attach(r http.Handler, opts ...attachOption) http.Handler {
	return newAttachHandler(e, r, opts...)
}

func (e *Engine) ServeHTTP(W http.ResponseWriter, r *http.Request) {

}
