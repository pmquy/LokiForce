package mq

import (
	"context"
	"sync"
)

type MessageQueue interface {
	Publish(ctx context.Context, topic string, payload any) error
	Subscribe(topic string, handler func(payload any))
}

type inMemoryMQ struct {
	mu          sync.RWMutex
	subscribers map[string][]func(payload any)
}

func NewInMemoryMQ() MessageQueue {
	return &inMemoryMQ{
		subscribers: make(map[string][]func(payload any)),
	}
}

func (mq *inMemoryMQ) Publish(ctx context.Context, topic string, payload any) error {
	mq.mu.RLock()
	handlers, exists := mq.subscribers[topic]
	mq.mu.RUnlock()

	if !exists {
		return nil
	}

	for _, handler := range handlers {
		h := handler
		go func() {
			h(payload)
		}()
	}
	return nil
}

func (mq *inMemoryMQ) Subscribe(topic string, handler func(payload any)) {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	mq.subscribers[topic] = append(mq.subscribers[topic], handler)
}
