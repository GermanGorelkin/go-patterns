package runner

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRunner_withoutInterrupt_withoutTimeout(t *testing.T) {
	var isCompleteTask1, isCompleteTask2, isCompleteTask3 bool
	task := func(id int) {
		switch id {
		case 0:
			isCompleteTask1 = true
		case 1:
			isCompleteTask2 = true
		case 2:
			isCompleteTask3 = true
		}
	}
	r := New(10 * time.Minute)
	r.Add(task, task, task)

	err := r.Start()

	assert.Nil(t, err)
	assert.True(t, isCompleteTask1)
	assert.True(t, isCompleteTask2)
	assert.True(t, isCompleteTask3)
}

func TestRunner_withTimeout(t *testing.T) {
	var isCompleteTask1, isCompleteTask2, isCompleteTask3 bool
	task := func(id int) {
		time.Sleep(200 * time.Millisecond)
		switch id {
		case 0:
			isCompleteTask1 = true
		case 1:
			isCompleteTask2 = true
		case 2:
			isCompleteTask3 = true
		}
	}
	r := New(500 * time.Millisecond)
	r.Add(task, task, task)

	err := r.Start()

	assert.Equal(t, err, ErrTimeout)
	assert.True(t, isCompleteTask1)
	assert.True(t, isCompleteTask2)
	assert.False(t, isCompleteTask3)
}
