package latch_test

import (
	"testing"
	"time"

	"github.com/KyKyPy3/clean/pkg/latch"
	"github.com/stretchr/testify/require"
)

func TestCountdownLatchOK(t *testing.T) {
	latch := latch.NewCountDownLatch()
	latch.Add(2)
	go func() {
		time.Sleep(100 * time.Millisecond)
		latch.Done()
	}()
	go func() {
		time.Sleep(100 * time.Millisecond)
		latch.Done()
	}()

	timeout := latch.WaitWithTimeout(time.Second)
	require.False(t, timeout)
}

func TestCountdownLatchClose(t *testing.T) {
	latch := latch.NewCountDownLatch()
	latch.Add(2)
	go func() {
		latch.Close()
	}()
	timeout := latch.WaitWithTimeout(time.Second)
	require.False(t, timeout)
}

func TestCountdownLatchWithTimeout(t *testing.T) {
	latch := latch.NewCountDownLatch()
	latch.Add(2)

	timeout := latch.WaitWithTimeout(time.Second)
	require.True(t, timeout)
}

func TestCountdownLatchWithTimeoutNoAdd(t *testing.T) {
	latch := latch.NewCountDownLatch()

	timeout := latch.WaitWithTimeout(time.Second)
	require.False(t, timeout)
}
