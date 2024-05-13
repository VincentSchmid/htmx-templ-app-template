package events

import (
	"sync"
)

type EventManager interface {
	OnEvent(event Event, handler func(interface{}))
	Emit(event Event)
}

type EventManagerImpl struct {
	listeners map[string][]chan interface{}
	events    []string
	lock      sync.Mutex
}

var _ EventManager = (*EventManagerImpl)(nil)

func NewEventManager() EventManager {
	return &EventManagerImpl{
		listeners: make(map[string][]chan interface{}),
		events:    make([]string, 0),
	}
}

type Event interface {
	Name() string
	Data() interface{}
}

func (em *EventManagerImpl) OnEvent(event Event, handler func(interface{})) {
	em.lock.Lock()
	defer em.lock.Unlock()

	ch := make(chan interface{})

	em.events = append(em.events, event.Name())
	em.listeners[event.Name()] = append(em.listeners[event.Name()], ch)

	go func() {
		for data := range ch {
			handler(data)
		}
	}()
}

func (em *EventManagerImpl) Emit(event Event) {
	em.lock.Lock()
	defer em.lock.Unlock()

	if listeners, ok := em.listeners[event.Name()]; ok {
		for _, listener := range listeners {
			go func(listener chan interface{}) {
				listener <- event.Data()
			}(listener)
		}
	}
}
