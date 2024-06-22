package engine

type Option func(engine *Engine)

type Engine struct {
}

func New(opts ...Option) *Engine {
	engine := &Engine{}
	for _, opt := range opts {
		opt(engine)
	}
	return engine
}
