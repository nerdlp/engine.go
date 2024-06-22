package engine

import (
	"context"
	"errors"
	"log/slog"
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

	///// INTERNAL VARIABLES ///////////
	sessionsPools map[string]transportClient
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

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		e.get(w, r)
		return
	}
}

func (e *Engine) get(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has(qk_sid) {
		// If no sid in query, it should be a hand shake
		request, err := e.prepareHandshakeRequest(r)
		if err != nil {
			slog.Error("fail to prepare hand shake", err)
			handleError(w, err)
			return
		}

		response, err := e.handleHandshake(r.Context(), request)
		if err != nil {
			handleError(w, err)
			return
		}

		response.render(w)
		return
	}

	// If there is sid in query, it should be a data polling from client
	return
}

// Send will add to the
func (e *Engine) Send(ctx context.Context, request *sendMessageRequest) error {
	session, exist := e.sessionsPools[request.sid]
	if !exist {
		return errors.New("session not found")
	}

	return session.Send(ctx, request)
}
