package core

import "context"

type CommandHandler[C any] interface {
	Handle(ctx context.Context, command C) error
}
