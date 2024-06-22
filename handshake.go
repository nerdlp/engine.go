package engine

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"time"
)

type handshakeRequest struct {
	eio       int32
	transport Transport
}

type handshakeResponse struct {
	Sid          string        `json:"sid"`
	Upgrades     []Transport   `json:"upgrades"`
	PingInterval time.Duration `json:"pingInterval"`
	PingTimeout  time.Duration `json:"pingTimeout"`
	MaxPayload   uint32        `json:"maxPayload"`
}

func (e *Engine) handleHandshake(ctx context.Context, request *handshakeRequest) (*handshakeResponse, error) {
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
		Sid:          sid.String(),
		Upgrades:     e.transports,
		PingInterval: e.pingInterval,
		PingTimeout:  e.pingTimeout,
		MaxPayload:   e.maxPayload,
	}, nil
}
