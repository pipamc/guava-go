package retry

import (
	"github.com/pipamc/guava-go/retry/strategy"
	"context"
	"sync/atomic"
)

type Callable func () error

type Retry struct {
	call Callable
	strategies []strategy.Strategy
}

func NewRetry(call Callable, strategies []strategy.Strategy) *Retry {
	return &Retry{call: call, strategies: strategies}
}

func (r *Retry) AppendStrategies(strategies ...strategy.Strategy) *Retry {
	r.strategies = append(r.strategies, strategies...)
	return r
}

func (r *Retry) SetRetryTime(rt uint) *Retry {
	return r
}

func (r *Retry) Call(ctx context.Context) error {
	if len(r.strategies) == 0 {
		return r.call()
	}

	var (
		err error
		interrupt uint32
	)

	finish := make(chan struct{})

	go func() {
		for attempt := uint(0); (attempt == 0 || err != nil) && r.shouldRetry(attempt, err) && !atomic.CompareAndSwapUint32(&interrupt, 1, 0); attempt++ {
			err = r.call()
		}
		close(finish)
	}()

	select {
	case <-ctx.Done():
		return nil
	case <- finish:
		return err
	}
}

func (r *Retry) shouldRetry(retryTime uint, err error) bool {
	shouldRetry := true
	for i, repeat := 0, len(r.strategies); shouldRetry && i < repeat; i++ {
		shouldRetry = shouldRetry && r.strategies[i](retryTime)
	}
	return shouldRetry
}