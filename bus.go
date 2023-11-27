package pubsub

import (
	"github.com/google/uuid"
	"reflect"
	"sync"
)

type EventHandler = func(args ...interface{})

type Bus struct {
	handlers map[string][]busHandler
	lock     *sync.Mutex
	wg       *sync.WaitGroup
}

type busHandler struct {
	fn       reflect.Value
	uniqueId string
}

func New() *Bus {
	return &Bus{
		handlers: make(map[string][]busHandler),
		lock:     &sync.Mutex{},
		wg:       &sync.WaitGroup{},
	}
}

func (b *Bus) Subscribe(topic string, handler EventHandler) string {
	b.lock.Lock()
	defer b.lock.Unlock()
	id := uuid.NewString()
	b.handlers[topic] = append(b.handlers[topic], busHandler{
		fn:       reflect.ValueOf(handler),
		uniqueId: id,
	})
	return id
}

func (b *Bus) Unsubscribe(topic string, id string) {
	handlers := b.handlers[topic]
	newHandlers := make([]busHandler, 0, len(handlers))
	for _, handler := range handlers {
		if handler.uniqueId == id {
			continue
		}
		newHandlers = append(newHandlers, handler)
	}
	b.lock.Lock()
	defer b.lock.Unlock()
	b.handlers[topic] = newHandlers
	return
}

func (b *Bus) Publish(topic string, args ...interface{}) {
	handlers := b.handlers[topic]
	for _, handler := range handlers {
		temp := handler
		b.wg.Add(1)
		go func() {
			defer b.wg.Done()
			callHandler(temp, args...)
		}()
	}
}

func (b *Bus) WaitUntilDone() {
	b.wg.Wait()
}

func callHandler(handler busHandler, args ...interface{}) {
	handlerArgsCount := handler.fn.Type().NumIn()
	if len(args) < handlerArgsCount {
		return
	}
	calledArgs := make([]reflect.Value, 0)
	for i := 0; i < handlerArgsCount; i++ {
		calledArgs = append(calledArgs, reflect.ValueOf(args[i]))
	}
	handler.fn.Call(calledArgs)
}
