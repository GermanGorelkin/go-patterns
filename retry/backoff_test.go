package retry

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
	"time"
)

func TestBackoff(t *testing.T) {
	factor := 100 * time.Millisecond
	maxAttempt := 5
	backoff := NewBackoff(100*time.Millisecond, 10*time.Second, maxAttempt, ConstantBackoff(factor))

	for i := 0; i < maxAttempt; i++ {
		assert.Equal(t, factor, backoff.Next())
	}
	assert.Equal(t, Stop, backoff.Next())

	backoff.Reset()

	for i := 0; i < maxAttempt; i++ {
		assert.Equal(t, factor, backoff.Next())
	}
	assert.Equal(t, Stop, backoff.Next())
}

func TestExponentialBackoff(t *testing.T) {
	tests := map[string]struct {
		attemptNum int
		min, max   time.Duration

		wantMin, wantMax time.Duration
	}{
		"1": {
			attemptNum: 1,
			min:        100 * time.Millisecond,
			max:        10 * time.Second,
			wantMin:    time.Duration(math.Pow(2, float64(1))),
			wantMax:    10 * time.Second,
		},
		"2": {
			attemptNum: 2,
			min:        100 * time.Millisecond,
			max:        10 * time.Second,
			wantMin:    time.Duration(math.Pow(2, float64(2))),
			wantMax:    10 * time.Second,
		},
		"3": {
			attemptNum: 3,
			min:        100 * time.Millisecond,
			max:        10 * time.Second,
			wantMin:    time.Duration(math.Pow(2, float64(3))),
			wantMax:    10 * time.Second,
		},
		"4": {
			attemptNum: 4,
			min:        100 * time.Millisecond,
			max:        10 * time.Second,
			wantMin:    time.Duration(math.Pow(2, float64(4))),
			wantMax:    10 * time.Second,
		},
		"over": {
			attemptNum: 11,
			min:        100 * time.Millisecond,
			max:        10 * time.Second,
			wantMin:    10 * time.Second,
			wantMax:    10 * time.Second,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := ExponentialBackoff(tc.attemptNum, tc.min, tc.max)
			between(t, actual, tc.wantMin, tc.wantMax)
		})
	}
}

func TestLinerBackoff(t *testing.T) {
	factor := 100 * time.Millisecond
	backoff := LinerBackoff(factor)

	tests := map[string]struct {
		attemptNum int
		min, max   time.Duration

		wantMin, wantMax time.Duration
	}{
		"1": {
			attemptNum: 1,
			min:        100 * time.Millisecond,
			max:        10 * time.Second,
			wantMin:    factor * 1,
			wantMax:    10 * time.Second,
		},
		"2": {
			attemptNum: 2,
			min:        100 * time.Millisecond,
			max:        10 * time.Second,
			wantMin:    factor * 2,
			wantMax:    10 * time.Second,
		},
		"3": {
			attemptNum: 3,
			min:        100 * time.Millisecond,
			max:        10 * time.Second,
			wantMin:    factor * 3,
			wantMax:    10 * time.Second,
		},
		"4": {
			attemptNum: 4,
			min:        100 * time.Millisecond,
			max:        10 * time.Second,
			wantMin:    factor * 4,
			wantMax:    10 * time.Second,
		},
		"over": {
			attemptNum: 111111111,
			min:        100 * time.Millisecond,
			max:        10 * time.Second,
			wantMin:    10 * time.Second,
			wantMax:    10 * time.Second,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := backoff(tc.attemptNum, tc.min, tc.max)
			between(t, actual, tc.wantMin, tc.wantMax)
		})
	}

}

func TestConstantBackoff(t *testing.T) {
	factor := 100 * time.Millisecond
	backoff := ConstantBackoff(factor)

	tests := map[string]struct {
		attemptNum int
		min, max   time.Duration
	}{
		"1": {
			attemptNum: 1,
			min:        100 * time.Millisecond,
			max:        10 * time.Second,
		},
		"2": {
			attemptNum: 2,
			min:        100 * time.Millisecond,
			max:        10 * time.Second,
		},
		"3": {
			attemptNum: 3,
			min:        100 * time.Millisecond,
			max:        10 * time.Second,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := backoff(tc.attemptNum, tc.min, tc.max)
			if actual != factor {
				t.Fatalf("Got %s, Expecting %s", actual, factor)
			}
		})
	}
}

func between(t *testing.T, actual, low, high time.Duration) {
	t.Helper()
	if actual < low {
		t.Fatalf("Got %s, Expecting >= %s", actual, low)
	}
	if actual > high {
		t.Fatalf("Got %s, Expecting <= %s", actual, high)
	}
}
