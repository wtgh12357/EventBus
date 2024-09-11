package main

import "sync"

type Event struct {
	Payload any
}

type EventChan chan Event

type EventBus struct {
	mu          sync.RWMutex
	subscribers map[string][]EventChan
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]EventChan),
	}
}

func (eb *EventBus) Publish(topic string, payload any) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	// 发布事件，通过事件名找到对应事件列表，遍历发送事件
	for _, subscriber := range eb.subscribers[topic] {
		subscriber <- Event{Payload: payload}
	}
}

func (eb *EventBus) Subscribe(event string) EventChan {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	ch := make(EventChan)

	// 向对应事件名列表中追加新的事件通道
	eb.subscribers[event] = append(eb.subscribers[event], ch)
	return ch
}

func (eb *EventBus) Unsubscribe(topic string, ch EventChan) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if subs, ok := eb.subscribers[topic]; ok {
		for i, sub := range subs {
			if ch == sub {
				eb.subscribers[topic] = append(subs[:i], subs[i+1:]...)
				close(ch)
				// 将通道值全部取出，清空
				for range ch {
				}
				return
			}
		}
	}
}
