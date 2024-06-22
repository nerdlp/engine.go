package engine

import (
	"context"
	"io"
	"net/http"
	"strconv"
)

type postPollingRequest struct {
	eio       int32
	transport Transport
	sid       string
	data      []byte
}

type postPollingResponse struct {
	request *postPollingRequest
	data    []*packet
}

func (e *Engine) preparePostPollingRequest(r *http.Request) (*postPollingRequest, error) {
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

	responseBody, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, socketError{
			code:       http.StatusBadRequest,
			innerError: err,
			message:    "fail to read request body",
		}
	}

	return &postPollingRequest{
		eio:       int32(eioNumber),
		transport: Transport(transport),
		sid:       sid,
		data:      responseBody,
	}, nil
}

func (e *Engine) handlePostPolling(_ context.Context, request *postPollingRequest) (*postPollingResponse, error) {
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

	packets := pollingSession.decompressPacket(request.data)
	for _, packet := range packets {
		pollingSession.receivePacket(packet)
	}

	return &postPollingResponse{
		data:    packets,
		request: request,
	}, nil
}

func (r *postPollingResponse) render(w http.ResponseWriter) {
	// case polling
	w.Header().Add(contentType, contentTypeOctecStream)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
