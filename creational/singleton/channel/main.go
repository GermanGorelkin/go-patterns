package channel

var addCh chan bool = make(chan bool)
var getCountCh chan chan int = make(chan chan int)
var quitCh chan bool = make(chan bool)

//func init() {
//	var count int
//
//	go func(addCh <-chan bool, getCountCh <-chan chan int, quitCh <-chan bool) {
//		for {
//			select {
//			case <-addCh:
//				count++
//			case ch := <-getCountCh:
//				ch <- count
//			case <-quitCh:
//				return
//			}
//		}
//	}(addCh, getCountCh, quitCh)
//}

type Singleton interface {
	AddOne()
	GetCount() int
	Stop()
}

type singleton struct {
	count int
}

var instance *singleton

func GetInstance() Singleton {
	if instance==nil {
		instance = new(singleton)

		go func(addCh <-chan bool, getCountCh <-chan chan int, quitCh <-chan bool) {
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
		}(addCh, getCountCh, quitCh)
	}

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
}
