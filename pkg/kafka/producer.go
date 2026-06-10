package kafka

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
)

type EventProducer interface {
	Publish(ctx context.Context, topic string, key, value []byte) error
}

// Producer — тонкая обёртка над kafka.Writer.
// Topic задаётся на каждое сообщение, поэтому один Producer обслуживает любые топики.
type Producer struct {
	writer *kafka.Writer
	log    *slog.Logger
}

func NewProducer(brokers []string, log *slog.Logger) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Balancer:     &kafka.Hash{},
			RequiredAcks: kafka.RequireAll,
			Async:        false,
		},
		log: log,
	}
}

// EnsureTopics создаёт топики при старте приложения, если их ещё нет.
// Admin-запрос автоматически маршрутизируется на контроллер кластера.
// Уже существующие топики не считаются ошибкой (идемпотентно).
func (p *Producer) EnsureTopics(ctx context.Context, topics ...kafka.TopicConfig) error {
	client := &kafka.Client{Addr: p.writer.Addr}
	resp, err := client.CreateTopics(ctx, &kafka.CreateTopicsRequest{Topics: topics})
	if err != nil {
		return fmt.Errorf("create topics: %w", err)
	}
	for topic, topicErr := range resp.Errors {
		if topicErr != nil && !errors.Is(topicErr, kafka.TopicAlreadyExists) {
			return fmt.Errorf("create topic %q: %w", topic, topicErr)
		}
		p.log.Info("kafka topic ready", slog.String("topic", topic))
	}
	return nil
}

// Publish отправляет одно сообщение синхронно. key используется для партиционирования
// (сообщения с одним key попадают в одну партицию и сохраняют порядок).
func (p *Producer) Publish(ctx context.Context, topic string, key, value []byte) error {
	if err := p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}); err != nil {
		p.log.Error("kafka publish failed", slog.String("topic", topic), slog.String("err", err.Error()))
		return err
	}
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
