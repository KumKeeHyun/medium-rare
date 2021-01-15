// Package erouter ...
// event-router
package erouter

import (
	"context"
	"sync"

	"github.com/KumKeeHyun/medium-rare/trend-service/config"
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
)

type EventHandlerFunc func(key, value []byte)

type handlerEntry struct {
	handler EventHandlerFunc
	ch      chan *sarama.ConsumerMessage
}

type eventRouter struct {
	routeTable map[string]*handlerEntry
	group      string
	ctx        context.Context
	ctxCancel  context.CancelFunc
	wg         sync.WaitGroup
	ready      chan bool
	log        *zap.Logger
}

func NewEventRouter(grp string, log *zap.Logger) *eventRouter {
	ctx, cancel := context.WithCancel(context.Background())
	return &eventRouter{
		routeTable: map[string]*handlerEntry{},
		group:      grp,
		ctx:        ctx,
		ctxCancel:  cancel,
		wg:         sync.WaitGroup{},
		ready:      make(chan bool),
		log:        log,
	}
}

func (er *eventRouter) SetHandler(topic string, handler EventHandlerFunc) {
	if _, exists := er.routeTable[topic]; exists {
		er.log.Panic("Topic conflict",
			zap.String("topic", topic))
	}

	er.routeTable[topic] = &handlerEntry{
		handler: handler,
		ch:      make(chan *sarama.ConsumerMessage, 100),
	}
}

func (er *eventRouter) StartRouter() error {
	er.log.Info("Start event router")

	for topic, entry := range er.routeTable {
		er.wg.Add(1)
		go func(t string, e *handlerEntry) {
			defer er.wg.Done()
			for {
				select {
				case msg := <-e.ch:
					e.handler(msg.Key, msg.Value)
				case <-er.ctx.Done():
					er.log.Info("Stop event handler",
						zap.String("topic", t))
					return
				}
			}
		}(topic, entry)
	}
	er.log.Info("Event handlers start running",
		zap.Int("handlers", len(er.routeTable)))

	if err := er.startConsumerGroup(); err != nil {
		er.log.Error("Error starting consumer group")
		return err
	}

	return nil
}

func (er *eventRouter) Stop() {
	er.ctxCancel()
	er.wg.Wait()
}

func (er *eventRouter) startConsumerGroup() error {
	// temp broker
	config.App.KafkaConfig.Brokers = []string{"220.70.2.5:8082"}

	er.log.Info("Create and start sarama consumer group",
		zap.Strings("brokers", config.App.KafkaConfig.Brokers))

	topics := make([]string, 0)
	for topic := range er.routeTable {
		topics = append(topics, topic)
	}

	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_5_0_0
	cfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest

	client, err := sarama.NewConsumerGroup(config.App.KafkaConfig.Brokers, er.group, cfg)
	if err != nil {
		er.log.Error("Error creating consumer group client",
			zap.Error(err))
		return err
	}

	er.wg.Add(1)
	go func() {
		defer er.wg.Done()
		for {
			if err := client.Consume(er.ctx, topics, er); err != nil {
				er.log.Panic("Error from sarama client",
					zap.Error(err))
			}

			if er.ctx.Err() != nil {
				er.log.Info("Stop consumer group client",
					zap.Error(er.ctx.Err()))
				return
			}
			er.ready = make(chan bool)
		}
	}()
	<-er.ready
	er.log.Info("Sarama consumer group start running")

	return nil
}

func (er *eventRouter) Setup(sarama.ConsumerGroupSession) error {
	close(er.ready)
	return nil
}

func (er *eventRouter) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (er *eventRouter) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		entry, exists := er.routeTable[msg.Topic]
		if !exists {
			er.log.Debug("Topic not initialized",
				zap.String("topic", msg.Topic))
			continue
		}
		entry.ch <- msg
	}

	return nil
}
