package usecase

import "context"

type Ping struct{}

func NewPing() *Ping {
	return &Ping{}
}

// Execute performs a trivial liveness check, always returning "pong".
func (u *Ping) Execute(context.Context) string {
	return "pong"
}
