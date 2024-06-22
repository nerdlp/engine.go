package engine

import (
	"context"
	"net/http"
	"strconv"
)

type getPollingRequest struct {
	eio       int32
	transport Transport
	sid       string
}

type getPollingResponse struct {
	request *getPollingRequest
	data    []byte
}

func (e *Engine) preparegetPollingRequest(r *http.Request) (*getPollingRequest, error) {
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
	if transport != Polling {
		return nil, socketError{
			code:    http.StatusBadRequest,
			message: "invalid transport method",
		}
	}

	// check sid
	sid := query.Get(qk_sid)
	if len(sid) == 0 {
		return nil, socketError{
			code:    http.StatusBadRequest,
			message: "invalid sid",
		}
	}

	return &getPollingRequest{
		eio:       int32(eioNumber),
		transport: Transport(transport),
		sid:       sid,
	}, nil
}

func (e *Engine) handleGetPolling(_ context.Context, request *getPollingRequest) (*getPollingResponse, error) {
	session, ok := e.sessionsPools[request.sid]
	if !ok {
		return nil, &socketError{
			code:    http.StatusBadRequest,
			message: "invalid sid",
		}
	}

	pollingSession, ok := session.(*pollingClient)
	if !ok {
		return nil, &socketError{
			code:    http.StatusBadRequest,
			message: "session is not polling",
		}
	}

	compressedPacket := pollingSession.compressPacket()

	return &getPollingResponse{
		data:    compressedPacket,
		request: request,
	}, nil
}

func (r *getPollingResponse) render(w http.ResponseWriter) {
	// case polling
	w.Header().Add(contentType, contentTypeOctecStream)
	w.WriteHeader(http.StatusOK)
	w.Write(r.data)
}
