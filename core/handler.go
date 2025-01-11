package core

import (
	"context"
	"github.com/patyhank/ptd/core/event"
	"github.com/sirupsen/logrus"
	"runtime/debug"
	"sync"
)

type EventManager struct {
	eventListenerMutex sync.Mutex
	eventListeners     []EventListener
	asyncCall          bool
	logger             logrus.Logger
}

type EventListener interface {
	OnEvent(client *EventClient, event event.Event)
}

type listenerFunc[E event.Event] struct {
	f func(client *EventClient, e E)
}

func (l *listenerFunc[E]) OnEvent(c *EventClient, e event.Event) {
	if evt, ok := e.(E); ok {
		l.f(c, evt)
	}
}

func NewListenerFunc[E event.Event](f func(client *EventClient, e E)) EventListener {
	return &listenerFunc[E]{f: f}
}

func (e *EventManager) DispatchEvent(client *Client, event event.Event) {
	defer func() {
		if r := recover(); r != nil {
			e.logger.WithFields(map[string]interface{}{
				"arg": r,
			}).Errorf("recovered panic in event listener\n%s", string(debug.Stack()))
			return
		}
	}()
	e.eventListenerMutex.Lock()
	defer e.eventListenerMutex.Unlock()
	c := &EventClient{Client: client, name: event.Name()}
	for i := range e.eventListeners {
		go func(i int) {
			defer func() {
				if r := recover(); r != nil {
					e.logger.WithFields(map[string]interface{}{
						"arg": r,
					}).Errorf("recovered panic in event listener\n%s", string(debug.Stack()))
					return
				}
			}()
			e.eventListeners[i].OnEvent(c, event)
		}(i)
	}
}

func (e *EventManager) AddEventListeners(listeners ...EventListener) {
	e.eventListenerMutex.Lock()
	defer e.eventListenerMutex.Unlock()
	e.eventListeners = append(e.eventListeners, listeners...)
}

func (e *EventManager) RemoveEventListeners(listeners ...EventListener) {
	e.eventListenerMutex.Lock()
	defer e.eventListenerMutex.Unlock()
	for _, listener := range listeners {
		for i, l := range e.eventListeners {
			if l == listener {
				e.eventListeners = append(e.eventListeners[:i], e.eventListeners[i+1:]...)
				break
			}
		}
	}
}

type EventClient struct {
	*Client
	name string
}

func (g *EventClient) Wait(ctx context.Context) {
	g.Client.WaitUpdate(ctx, g.name)
}

func (g *EventClient) PrepareWait() {
	g.Client.PrepareWait(g.name)
}
