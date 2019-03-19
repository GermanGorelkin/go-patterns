package classic

type Observer interface {
	Update()
}

type Subject interface {
	Attach(Observer)
	Detach(Observer)
	//Notify()
}

type ConcreteSubject struct {
	observers []Observer
	state     interface{}
}

func (s *ConcreteSubject) Attach(o Observer) {
	s.observers = append(s.observers, o)
}
func (s *ConcreteSubject) Detach(o Observer) {
	for i, v := range s.observers {
		if o == v {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			break
		}
	}
}
func (s *ConcreteSubject) notify() {
	for _, v := range s.observers {
		v.Update()
	}
}
func (s *ConcreteSubject) SetState(st interface{}) {
	s.state = st
	s.notify()
}
func (s *ConcreteSubject) GetState() interface{} {
	return s.state
}

type ConcreteObserver struct {
	subject *ConcreteSubject
	state   interface{}
}

func (o *ConcreteObserver) Update() {
	o.state = o.subject.GetState()
}
