package events

import (
	"fmt"
	"log/slog"
	"path/filepath"
)

type EventHandler func(Event) error

type EventBus interface {
	Emit(...Event) error
	Listen(string, EventHandler) error
}

type InMemoryEventBus struct {
	logger    *slog.Logger
	listeners map[string][]EventHandler
	//eventsQueue chan Event
}

func NewEventBus(logger *slog.Logger) EventBus {
	return &InMemoryEventBus{
		logger:    logger,
		listeners: make(map[string][]EventHandler),
	}
}

func (b *InMemoryEventBus) Emit(events ...Event) error {
	for _, event := range events {
		for pattern, listeners := range b.listeners {
			match, err := filepath.Match(pattern, event.Name)
			if !match || err != nil {
				continue
			}

			for _, l := range listeners {
				go func(e Event) {
					defer b.recoverFunction()
					err := l(e)
					if err != nil {
						b.logger.Error(err.Error())
					}
				}(event)
			}
		}
	}
	return nil
}

func (b *InMemoryEventBus) Listen(pattern string, handler EventHandler) error {
	listeners := b.listeners[pattern]
	if listeners == nil {
		listeners = make([]EventHandler, 0)
	}

	listeners = append(listeners, handler)
	b.listeners[pattern] = listeners

	return nil
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
