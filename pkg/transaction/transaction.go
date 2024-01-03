package transaction

import (
	"context"
	"errors"
	"fmt"
)

var ErrTxNotFound = errors.New("no transaction found in context")

type txKey struct{}

var txID txKey

type TxFunc[T Tx] func(context.Context, T) error

type Tx interface {
	Rollback() error
	Commit() error
}

func setTxToContext(ctx context.Context, tx Tx) context.Context {
	return context.WithValue(ctx, txID, tx)
}

func GetTxFromContext[T Tx](ctx context.Context) (T, error) {
	tx, ok := ctx.Value(txID).(T)
	if !ok {
		var zero T
		return zero, ErrTxNotFound
	}
	return tx, nil
}

type Transaction[T Tx] struct {
	txFactory func(context.Context) (Tx, error)
}

func New[T Tx](txFactory func(context.Context) (Tx, error)) *Transaction[T] {
	return &Transaction[T]{
		txFactory: txFactory,
	}
}

func (tm *Transaction[T]) Current(ctx context.Context, fn TxFunc[T]) error {
	t := ctx.Value(txID)
	if t == nil {
		return tm.makeTxHandler(ctx, fn)
	}

	tx, ok := t.(T)
	if !ok {
		return ErrTxNotFound
	}

	return fn(ctx, tx)
}

func (tm *Transaction[T]) New(ctx context.Context, fn TxFunc[T]) error {
	return tm.makeTxHandler(ctx, fn)
}

func (tm *Transaction[T]) makeTxHandler(ctx context.Context, fn TxFunc[T]) error {
	// Begin Transaction
	t, err := tm.txFactory(ctx)
	if err != nil {
		return err
	}
	tx, ok := t.(T)
	if !ok {
		var zero T
		return fmt.Errorf("factory transaction produced wrong type: want %T, got %T", zero, t)
	}
	c := setTxToContext(ctx, tx)
	defer func() {
		_ = tx.Rollback()
	}()

	err = fn(c, tx)
	if err != nil {
		return err
	}

	_ = tx.Commit()
	return nil
}
