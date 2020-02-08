package main

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestCounter(t *testing.T) {
	c := new(counter)

	var wg sync.WaitGroup
	numLoop := 1000

	wg.Add(numLoop)
	for i := 0; i < numLoop; i++ {
		go func() {
			defer wg.Done()
			c.Increment()
		}()
	}

	wg.Add(numLoop)
	for i := 0; i < numLoop; i++ {
		go func() {
			defer wg.Done()
			c.Decrement()
		}()
	}

	wg.Wait()

	expected := 0
	assert.Equal(t, expected, c.count)
}

func BenchmarkCountV1(b *testing.B) {
	c := new(counter)
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				c.CountV1()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkCountV2(b *testing.B) {
	c := new(counter)
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				c.CountV2()
			}()
		}
		wg.Wait()
	}
}
