// Package erouter ...
// event-router
package erouter

import (
	"context"
	"sync"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"
)

type EventHandlerFunc func(key, value []byte)

type ConsumerGroup struct {
	MsgOutput      chan *sarama.ConsumerMessage
	ConsumerCtx    context.Context
	ConsumerCancel context.CancelFunc
}

type eventRouter struct {
	consumerGroup *ConsumerGroup
	routeTable    map[string]EventHandlerFunc
	handlerGrp    sync.WaitGroup
	log           *zap.Logger
}

func NewEventRouter(cg *ConsumerGroup, log *zap.Logger) *eventRouter {
	return &eventRouter{
		consumerGroup: cg,
		routeTable:    map[string]EventHandlerFunc{},
		handlerGrp:    sync.WaitGroup{},
		log:           log,
	}
}

func (er *eventRouter) SetHandler(topic string, handler EventHandlerFunc) {
	if _, exists := er.routeTable[topic]; exists {
		er.log.Panic("Topic conflict",
			zap.String("topic", topic))
	}
	er.routeTable[topic] = handler
}

func (er *eventRouter) Start() {
	chanTable := map[string]chan *sarama.ConsumerMessage{}
	for topic := range er.routeTable {
		chanTable[topic] = make(chan *sarama.ConsumerMessage, 100)
	}

	for topic, handler := range er.routeTable {
		er.handlerGrp.Add(1)

		go func(t string, h EventHandlerFunc) {
			defer er.handlerGrp.Done()
			ch := chanTable[t]

			for {
				select {
				case msg := <-ch:
					h(msg.Key, msg.Value)
				case <-er.consumerGroup.ConsumerCtx.Done():
					er.log.Info("Stop event handler",
						zap.String("topic", t))
					return
				}
			}
		}(topic, handler)
	}

	er.handlerGrp.Add(1)
	go func() {
		defer er.handlerGrp.Done()

		for {
			select {
			case msg := <-er.consumerGroup.MsgOutput:
				if ch, exists := chanTable[msg.Topic]; exists {
					ch <- msg
				} else {
					er.log.Debug("Topic not initialized",
						zap.String("topic", msg.Topic))
				}
			case <-er.consumerGroup.ConsumerCtx.Done():
				er.log.Info("Stop event router")
				return
			}
		}
	}()
}

func (er *eventRouter) Stop() {
	er.consumerGroup.ConsumerCancel()
	er.handlerGrp.Wait()
}
