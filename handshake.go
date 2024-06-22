package engine

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type handshakeRequest struct {
	eio       int32
	transport Transport
}

type handshakeResponse struct {
	request *handshakeRequest
	data    *handshakeResponseData
}

type handshakeResponseData struct {
	Sid          string        `json:"sid"`
	Upgrades     []Transport   `json:"upgrades"`
	PingInterval time.Duration `json:"pingInterval"`
	PingTimeout  time.Duration `json:"pingTimeout"`
	MaxPayload   uint32        `json:"maxPayload"`
}

func (e *Engine) prepareHandshakeRequest(r *http.Request) (*handshakeRequest, error) {
	query := r.URL.Query()
	// check engine.io format
	eio := query.Get(qk_EIO)
	eioNumber, err := strconv.Atoi(eio)
	if err != nil {
		err = socketError{
			innerError: err,
			code:       http.StatusBadRequest,
			message:    "invalid eio format",
		}
		return nil, err
	}
	if eioNumber != 4 {
		err = socketError{
			innerError: err,
			code:       http.StatusBadRequest,
			message:    "invalid eio version",
		}
		return nil, err
	}

	// check transport format
	transport := Transport(query.Get(qk_Transport))
	if transport != Polling && transport != WebSocket {
		return nil, socketError{
			code:    http.StatusBadRequest,
			message: "invalid transport method",
		}
	}

	return &handshakeRequest{
		eio:       int32(eioNumber),
		transport: Transport(transport),
	}, nil
}

func (e *Engine) handleHandshake(_ context.Context, request *handshakeRequest) (*handshakeResponse, error) {
	sid, err := uuid.NewUUID()
	if err != nil {
		err = socketError{
			innerError: err,
			code:       http.StatusInternalServerError,
			message:    "create uuid failed",
		}
		slog.Error(err.Error())
		return nil, err
	}

	return &handshakeResponse{
		data: &handshakeResponseData{
			Sid:          sid.String(),
			Upgrades:     e.transports,
			PingInterval: e.pingInterval,
			PingTimeout:  e.pingTimeout,
			MaxPayload:   e.maxPayload,
		},
		request: request,
	}, nil
}

func (r *handshakeResponse) render(w http.ResponseWriter) {
	data, err := json.Marshal(r.data)
	if err != nil {
		handleError(w, err)
		return
	}
	packet := newPacket(open, data)

	switch r.request.transport {
	case Polling:
		// case polling
		w.Header().Add(contentType, contentTypeOctecStream)
		w.WriteHeader(http.StatusOK)
		w.Write(packet.encodePolling())
		return
	case WebSocket:
		// case websocket
		w.Header().Add(contentType, contentTypeOctecStream)
		w.WriteHeader(http.StatusSwitchingProtocols)
		w.Write(packet.encodePolling())
		return
	default:
		handleError(w, &socketError{
			code:    http.StatusBadRequest,
			message: "unsuported protocol",
		})
		return
	}
}
