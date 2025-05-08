package producer

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/mephirious/advanced-programming-2/inventory-service/internal/adapter/nats/dto"
	"github.com/mephirious/advanced-programming-2/inventory-service/internal/domain"
	"github.com/mephirious/advanced-programming-2/inventory-service/pkg/nats"

	pb "github.com/mephirious/advanced-programming-2/inventory-service/proto/events"
)

const PushTimeout = time.Second * 30

type InventoryEventProducer struct {
	natsClient *nats.Client
	subject    string
}

func NewInventoryEventProducer(
	natsClient *nats.Client,
	subject string,
) *InventoryEventProducer {
	return &InventoryEventProducer{
		natsClient: natsClient,
		subject:    subject,
	}
}

func (p *InventoryEventProducer) Push(ctx context.Context, event *domain.Product, eventType pb.InventoryEventType) error {
	_, cancel := context.WithTimeout(ctx, PushTimeout)
	defer cancel()

	pbEvent := dto.ToInventoryEvent(event, eventType)
	data, err := proto.Marshal(pbEvent)
	if err != nil {
		return fmt.Errorf("proto.Marshal: %w", err)
	}

	err = p.natsClient.Conn.Publish(p.subject, data)
	if err != nil {
		return fmt.Errorf("p.natsClient.Conn.Publish: %w", err)
	}
	log.Printf("Inventory event is pushed: %+v [%s]", event, eventType)

	return nil
}
