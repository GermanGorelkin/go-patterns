package retry

import (
	"context"
	"time"
)

type Action int

const (
	Succeed Action = iota
	Fail
	Retry
)

type Worker func(ctx context.Context) error
type RetryPolicy func(err error) Action

type Retrier struct {
	backoff     *Backoff
	retryPolicy RetryPolicy
}

func NewRetrier(backoff *Backoff, retryPolicy RetryPolicy) Retrier {
	if retryPolicy == nil {
		retryPolicy = DefaultRetryPolicy
	}

	return Retrier{
		backoff:     backoff,
		retryPolicy: retryPolicy,
	}
}

func (r Retrier) Run(ctx context.Context, work Worker) error {
	for {
		err := work(ctx)

		switch r.retryPolicy(err) {
		case Succeed, Fail:
			return err
		case Retry:
			// log.Println(err) // error logging
			var delay time.Duration
			if delay = r.backoff.Next(); delay == Stop {
				return err
			}
			timeout := time.After(delay)
			if err := r.sleep(ctx, timeout); err != nil {
				return err
			}
		}
	}
}

func (r *Retrier) sleep(ctx context.Context, t <-chan time.Time) error {
	select {
	case <-t:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func DefaultRetryPolicy(err error) Action {
	if err == nil {
		return Succeed
	}
	return Retry
}
