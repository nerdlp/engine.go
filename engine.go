package engine

import (
	"net/http"
	"strconv"
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
	query := r.URL.Query()
	eio := query.Get(qk_EIO)
	transport := query.Get(qk_Transport)
	response := new(response)

	eioNumber, err := strconv.Atoi(eio)
	if err != nil {
		err = socketError{
			innerError: err,
			code:       http.StatusBadRequest,
			message:    "invalid eio format",
		}
		handleError(w, err)
		return
	}
	if eioNumber != 4 {
		err = socketError{
			innerError: err,
			code:       http.StatusBadRequest,
			message:    "invalid eio version",
		}
		handleError(w, err)
		return
	}

	switch Transport(transport) {
	case Polling:
		response.status = http.StatusOK
	case WebSocket:
		response.status = http.StatusSwitchingProtocols
	default:
		err = socketError{
			code:    http.StatusBadRequest,
			message: "invalid transport method",
		}
		handleError(w, err)
		return
	}

	responseBody, err := e.handleHandshake(r.Context(), &handshakeRequest{
		eio:       int32(eioNumber),
		transport: Transport(transport),
	})
	if err != nil {
		handleError(w, err)
		return
	}
	response.body = responseBody

	handleResponse(w, response)
	return
}
