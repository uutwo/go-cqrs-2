package cqrs

import (
	"fmt"
	"github.com/satori/go.uuid"
)

type EventFactoryFunc func(uuid.UUID, int) Event

// EventFactory is the interface that an event store should implement.
// An event factory returns instances of an event given the event type as a string.
type EventFactory interface {
	GetEvent(string, uuid.UUID, int) Event
}

// CallbackEventFactory uses callback/delegate functions to instantiate event instances
// given the name of the event type as a string.
type CallbackEventFactory struct {
	eventFactories map[string]EventFactoryFunc
}

// NewCallbackEventFactory constructs a new CallbackEventFactory
func NewCallbackEventFactory() *CallbackEventFactory {
	return &CallbackEventFactory{
		eventFactories: make(map[string]EventFactoryFunc),
	}
}

// RegisterCallback registers a delegate that will return an event instance given
// an event type name as a string.
func (t *CallbackEventFactory) RegisterCallback(callback EventFactoryFunc) error {
	typeName := callback(uuid.NewV4(), 0).EventType()
	if _, ok := t.eventFactories[typeName]; ok {
		return fmt.Errorf("Factory callback/delegate already registered for type: \"%s\"", typeName)
	}
	t.eventFactories[typeName] = callback
	return nil
}

// GetEvent returns an event instance given an event type as a string.
//
// An appropriate delegate must be registered for the event type.
// If an appropriate delegate is not registered, the method will return nil.
func (t *CallbackEventFactory) GetEvent(typeName string, aggregateId uuid.UUID, version int) Event {
	if f, ok := t.eventFactories[typeName]; ok {
		return f(aggregateId, version)
	}
	return nil
}
