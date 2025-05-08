package producer

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/mephirious/advanced-programming-2/order-service/internal/adapter/nats/dto"
	"github.com/mephirious/advanced-programming-2/order-service/internal/domain"
	"github.com/mephirious/advanced-programming-2/order-service/pkg/nats"

	pb "github.com/mephirious/advanced-programming-2/order-service/proto/events"
)

const PushTimeout = time.Second * 30

type OrderEventProducer struct {
	natsClient *nats.Client
	subject    string
}

func NewOrderEventProducer(
	natsClient *nats.Client,
	subject string,
) *OrderEventProducer {
	return &OrderEventProducer{
		natsClient: natsClient,
		subject:    subject,
	}
}

func (p *OrderEventProducer) Push(ctx context.Context, event *domain.Order, eventType pb.OrderEventType) error {
	_, cancel := context.WithTimeout(ctx, PushTimeout)
	defer cancel()

	pbEvent := dto.ToOrderEvent(event, eventType)
	data, err := proto.Marshal(pbEvent)
	if err != nil {
		return fmt.Errorf("proto.Marshal: %w", err)
	}

	err = p.natsClient.Conn.Publish(p.subject, data)
	if err != nil {
		return fmt.Errorf("p.natsClient.Conn.Publish: %w", err)
	}
	log.Printf("Order event is pushed: %+v [%s]", event, eventType)

	return nil
}
