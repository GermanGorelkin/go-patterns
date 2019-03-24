package pubsub

type Subscriber interface {
	Notify(interface{})
	Close()
}

type Publisher interface {
	start()
	AddSubscriber() chan<- Subscriber
	RemoveSubscriber() chan<- Subscriber
	PublishMessage() chan<- interface{}
	Stop()
}
