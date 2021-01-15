package util

import (
	"context"

	"github.com/KumKeeHyun/medium-rare/trend-service/config"
	"github.com/KumKeeHyun/medium-rare/trend-service/util/erouter"
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
)

func BuildConsumerGroup(log *zap.Logger) (*erouter.ConsumerGroup, error) {
	// temp broker
	config.App.KafkaConfig.Brokers = []string{"220.70.2.5:8082"}
	topics := []string{"read-article", "create-user"}

	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_5_0_0
	cfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest

	output := make(chan *sarama.ConsumerMessage, 100)
	ctx, cancel := context.WithCancel(context.Background())
	cg := &erouter.ConsumerGroup{
		MsgOutput:      output,
		ConsumerCtx:    ctx,
		ConsumerCancel: cancel,
	}

	c := consumer{
		ready:   make(chan bool),
		message: output,
	}

	client, err := sarama.NewConsumerGroup(config.App.KafkaConfig.Brokers, "trend", cfg)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			if err := client.Consume(ctx, topics, &c); err != nil {
				log.Panic("Error creating consumer group client",
					zap.Error(err))
			}

			if ctx.Err() != nil {
				log.Info("Stop consumer group client",
					zap.Error(ctx.Err()))
				return
			}
			c.ready = make(chan bool)
		}
	}()
	<-c.ready

	return cg, nil
}

type consumer struct {
	ready   chan bool
	message chan *sarama.ConsumerMessage
}

func (c *consumer) Setup(sarama.ConsumerGroupSession) error {
	close(c.ready)
	return nil
}

func (c *consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		c.message <- msg
	}

	return nil
}
