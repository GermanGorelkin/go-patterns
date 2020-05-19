package retry

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRetrier(t *testing.T) {
	t.Run("10 fails", func(t *testing.T) {
		factor := 100 * time.Millisecond
		maxAttempt := 10
		backoff := NewBackoff(100*time.Millisecond, 10*time.Second, maxAttempt, ConstantBackoff(factor))
		r := NewRetrier(backoff, nil)

		var attempts int
		retrierErr := errors.New("error")
		var elapsed []time.Duration

		start := time.Now()
		err := r.Run(context.Background(), func(ctx context.Context) error {
			attempts++
			elapsed = append(elapsed, time.Since(start))
			return retrierErr
		})

		assert.Equal(t, maxAttempt+1, attempts)
		assert.Equal(t, retrierErr, err)

		begin := time.Duration(0)
		end := factor
		for _, got := range elapsed {
			between(t, got, begin, end)
			begin += factor
			end += factor
		}
	})

	t.Run("2 fails and succeed", func(t *testing.T) {
		factor := 100 * time.Millisecond
		maxAttempt := 10
		backoff := NewBackoff(100*time.Millisecond, 10*time.Second, maxAttempt, ConstantBackoff(factor))
		r := NewRetrier(backoff, nil)

		var attempts int
		retrierErr := errors.New("error")
		var elapsed []time.Duration

		start := time.Now()
		err := r.Run(context.Background(), func(ctx context.Context) error {
			attempts++
			elapsed = append(elapsed, time.Since(start))

			if attempts == 2 {
				return nil
			}
			return retrierErr
		})

		assert.Equal(t, 2, attempts)
		assert.Nil(t, err)

		begin := time.Duration(0)
		end := factor
		for _, got := range elapsed {
			between(t, got, begin, end)
			begin += factor
			end += factor
		}
	})
}
