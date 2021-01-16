package util

import (
	"github.com/KumKeeHyun/medium-rare/article-service/config"
	"github.com/Shopify/sarama"
)

func BuildSyncProducer() (sarama.SyncProducer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 10
	cfg.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(config.App.KafkaConfig.Brokers, cfg)
	if err != nil {
		return nil, err
	}

	return producer, nil
}

func BuildMockSyncProducer() (sarama.SyncProducer, error) {
	return &MockSyncProducer{}, nil
}

type MockSyncProducer struct {
}

func (msp *MockSyncProducer) SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	return 0, 0, nil
}

func (msp *MockSyncProducer) SendMessages(msgs []*sarama.ProducerMessage) error {
	return nil
}

func (msp *MockSyncProducer) Close() error {
	return nil
}
