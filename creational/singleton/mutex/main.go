package mutex

import "sync"

type Singleton interface {
	AddOne()
	GetCount() int
}

type singleton struct {
	count int
	sync.RWMutex
}

var instance *singleton
var once sync.Once

func GetInstance() Singleton {
	once.Do(func() {
		instance = new(singleton)
	})

	return instance
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
