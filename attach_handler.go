package engine

import (
	"net/http"
	"strings"
)

const (
	defaultPath = "engine.io"
)

type attachOption func(h *attachHandler)

type attachHandler struct {
	engine  *Engine
	handler http.Handler
	// Properties
	path string
}

func newAttachHandler(engine *Engine, handler http.Handler, opts ...attachOption) *attachHandler {
	h := &attachHandler{
		engine:  engine,
		handler: handler,
		path:    defaultPath,
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

func (h *attachHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(ComputePath(h.path), ComputePath(r.URL.Path)) {
		h.engine.ServeHTTP(w, r)
		return
	}
	h.handler.ServeHTTP(w, r)
}
