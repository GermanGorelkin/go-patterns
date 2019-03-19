package evlistener

import (
	"fmt"
	"io"
)

type FileInfo struct {
	Name string
}

type EventType string

type Observer interface {
	Update(EventType, FileInfo)
}
type Subject interface {
	Attach(EventType, Observer)
	Detach(EventType, Observer)
}

type EventManager struct {
	observers map[EventType][]Observer
}

func (em *EventManager) Attach(et EventType, o Observer) {
	em.observers[et] = append(em.observers[et], o)
}
func (em *EventManager) Detach(et EventType, o Observer) {
	for i, v := range em.observers[et] {
		if v == o {
			em.observers[et] = append(em.observers[et][:i], em.observers[et][i+1:]...)
			break
		}
	}
}
func (em *EventManager) Notify(et EventType, fi FileInfo) {
	for _, o := range em.observers[et] {
		o.Update(et, fi)
	}
}

const (
	EventTypeCreate EventType = "create"
	EventTypeRemove EventType = "remove"
	EventTypeModify EventType = "modify"
)

type Watcher struct {
	eventManager *EventManager
}

func (w *Watcher) OnCreate(fi FileInfo) {
	w.eventManager.Notify(EventTypeCreate, fi)
}
func (w *Watcher) OnRemove(fi FileInfo) {
	w.eventManager.Notify(EventTypeRemove, fi)
}
func (w *Watcher) OnModify(fi FileInfo) {
	w.eventManager.Notify(EventTypeModify, fi)
}

type LoggingListener struct {
	store io.Writer
}

func (l *LoggingListener) Update(et EventType, fi FileInfo) {
	log := fmt.Sprintf("%s|%s", et, fi.Name)
	l.store.Write([]byte(log))
}

type AlertListener struct {
	service io.Writer
}

func (l *AlertListener) Update(et EventType, fi FileInfo) {
	log := fmt.Sprintf("%s|%s", et, fi.Name)
	l.service.Write([]byte(log))
}
