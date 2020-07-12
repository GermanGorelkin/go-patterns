package circuit_breaker

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func succeed(cb *CircuitBreaker) error {
	_, err := cb.Execute(func() (interface{}, error) { return nil, nil })
	return err
}

func fail(cb *CircuitBreaker) error {
	msg := "fail"
	_, err := cb.Execute(func() (interface{}, error) { return nil, fmt.Errorf(msg) })
	if err.Error() == msg {
		return nil
	}
	return err
}

func pseudoSleep(cb *CircuitBreaker, period time.Duration) {
	if !cb.expiry.IsZero() {
		cb.expiry = cb.expiry.Add(-period)
	}
}

func TestCircuitBreaker(t *testing.T) {
	cb := NewCircuitBreaker(Config{
		Name:        "test circuit breaker",
		MaxRequests: 2,
		ReadyToTrip: func(counts Counts) bool {
			return counts.ConsecutiveFailures > 5
		},
	})

	for i := 0; i < 5; i++ {
		assert.Nil(t, fail(cb))
	}

	assert.Equal(t, StateClosed, cb.state)
	assert.Equal(t, Counts{5, 0, 5, 0, 5}, cb.counts)

	assert.Nil(t, succeed(cb))
	assert.Equal(t, StateClosed, cb.state)
	assert.Equal(t, Counts{6, 1, 5, 1, 0}, cb.counts)

	assert.Nil(t, fail(cb))
	assert.Equal(t, StateClosed, cb.state)
	assert.Equal(t, Counts{7, 1, 6, 0, 1}, cb.counts)

	// StateClosed to StateOpen
	for i := 0; i < 5; i++ {
		assert.Nil(t, fail(cb)) // 6 consecutive failures
	}

	assert.Equal(t, StateOpen, cb.state)
	assert.Equal(t, Counts{0, 0, 0, 0, 0}, cb.counts)
	assert.False(t, cb.expiry.IsZero())

	assert.Error(t, succeed(cb))
	assert.Error(t, fail(cb))
	assert.Equal(t, Counts{0, 0, 0, 0, 0}, cb.counts)

	pseudoSleep(cb, time.Duration(59)*time.Second)
	assert.Equal(t, StateOpen, cb.state)

	// StateOpen to StateHalfOpen
	pseudoSleep(cb, time.Duration(1)*time.Second) // over Timeout
	assert.Nil(t, succeed(cb))
	assert.Equal(t, StateHalfOpen, cb.state)
	assert.True(t, cb.expiry.IsZero())
	assert.Equal(t, Counts{1, 1, 0, 1, 0}, cb.counts)

	// StateHalfOpen to StateOpen
	assert.Nil(t, fail(cb))
	assert.Equal(t, StateOpen, cb.state)
	assert.False(t, cb.expiry.IsZero())
	assert.Equal(t, Counts{0, 0, 0, 0, 0}, cb.counts)

	// StateOpen to StateHalfOpen
	pseudoSleep(cb, time.Duration(60)*time.Second) // over Timeout
	assert.Nil(t, succeed(cb))
	assert.Equal(t, StateHalfOpen, cb.state)
	assert.True(t, cb.expiry.IsZero())
	assert.Equal(t, Counts{1, 1, 0, 1, 0}, cb.counts)

	// StateHalfOpen to StateClosed
	assert.Nil(t, succeed(cb)) // ConsecutiveSuccesses(2) >= MaxRequests(2)
	assert.Equal(t, StateClosed, cb.state)
	assert.Equal(t, Counts{0, 0, 0, 0, 0}, cb.counts)
	assert.True(t, cb.expiry.IsZero())
}
