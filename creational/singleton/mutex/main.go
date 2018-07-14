package mutex

import "sync"

type singleton struct {
	count int
	sync.RWMutex
}

var instance singleton

func GetInstance() *singleton {
	return &instance
}

func (s *singleton) AddOne() {
	s.Lock()
	defer s.Unlock()
	s.count++
}

func (s *singleton) GetCount()int {
	s.RLock()
	defer s.RUnlock()
	return s.count
}
