package main

import "sync"

type counter struct {
	sync.RWMutex
	count int
}

func (c *counter) Increment() {
	c.Lock()
	defer c.Unlock()
	c.count++
}
func (c *counter) Decrement() {
	c.Lock()
	defer c.Unlock()
	c.count--
}
func (c *counter) CountV1() int {
	c.Lock()
	defer c.Unlock()
	return c.count
}
func (c *counter) CountV2() int {
	c.RLock()
	defer c.RUnlock()
	return c.count
}
