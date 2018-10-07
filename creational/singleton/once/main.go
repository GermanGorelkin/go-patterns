package once

import "sync"

type Singleton interface {
	AddOne() int
}

type singleton struct {
	count int
}

var instance *singleton
var once sync.Once

func GetInstance() Singleton {
	once.Do(func() {
		instance = new(singleton)
	})

	return instance
}

func (s *singleton) AddOne() int {
	s.count++
	return s.count
}
