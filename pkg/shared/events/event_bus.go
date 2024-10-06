package events

import (
	"fmt"
	"log/slog"
)

type EventHandler[T any] func(Event[T]) error

type EventBus interface {
	Emit(...Event[any]) error
	Listen(EventName, EventHandler[any])
}

type InMemoryEventBus struct {
	logger    *slog.Logger
	listeners map[EventName][]EventHandler[any]
}

func NewEventBus(logger *slog.Logger) EventBus {
	return &InMemoryEventBus{
		logger:    logger,
		listeners: make(map[EventName][]EventHandler[any]),
	}
}

func (b *InMemoryEventBus) Emit(events ...Event[any]) error {
	for _, event := range events {
		listeners, hasListeners := b.listeners[event.Name]
		if !hasListeners {
			continue
		}

		for _, l := range listeners {
			go func(e Event[any]) {
				defer b.recoverFunction()
				err := l(e)
				if err != nil {
					b.logger.Error(err.Error())
				}
			}(event)
		}
	}
	return nil
}

func (b *InMemoryEventBus) Listen(name EventName, handler EventHandler[any]) {
	listeners := b.listeners[name]
	if listeners == nil {
		listeners = make([]EventHandler[any], 0)
	}

	listeners = append(listeners, handler)
	b.listeners[name] = listeners
}

func (b *InMemoryEventBus) recoverFunction() {
	r := recover()
	if r == nil {
		return
	}

	if err, isErr := r.(error); isErr {
		b.logger.Error(err.Error())
		return
	}

	b.logger.Error(fmt.Sprint(r))
}
