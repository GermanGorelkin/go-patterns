package channel

import "sync"

var addCh chan bool = make(chan bool)
var getCountCh chan chan int = make(chan chan int)
var quitCh chan bool = make(chan bool)

type Singleton interface {
	AddOne()
	GetCount() int
	Stop()
}

type singleton struct {
	count int
}

var instance *singleton
var once sync.Once

func GetInstance() Singleton {
	once.Do(func() {
		instance = new(singleton)

		go func() {
			for {
				select {
				case <-addCh:
					instance.count++
				case ch := <-getCountCh:
					ch <- instance.count
				case <-quitCh:
					return
				}
			}
		}()
	})

	return instance
}

func (s *singleton) AddOne() {
	addCh <- true
}

func (s *singleton) GetCount() int {
	resCh := make(chan int)
	defer close(resCh)
	getCountCh <- resCh
	return <-resCh
}

func (s *singleton) Stop() {
	quitCh <- true
	close(addCh)
	close(getCountCh)
	close(quitCh)
	instance = nil
}
