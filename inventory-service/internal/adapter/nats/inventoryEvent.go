package producer

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/mephirious/advanced-programming-2/inventory-service/internal/domain"
	"github.com/mephirious/advanced-programming-2/inventory-service/pkg/nats"
	pb "github.com/mephirious/advanced-programming-2/inventory-service/proto/events"
)

const PushTimeout = time.Second * 30

type InventoryEventProducer struct {
	natsClient *nats.Client
	subject    string
}

func NewInventoryEventProducer(natsClient *nats.Client, subject string) *InventoryEventProducer {
	return &InventoryEventProducer{
		natsClient: natsClient,
		subject:    subject,
	}
}

func (p *InventoryEventProducer) Push(ctx context.Context, event *domain.Product, eventType pb.InventoryEventType) error {
	ctx, cancel := context.WithTimeout(ctx, PushTimeout)
	defer cancel()

	if eventType == pb.InventoryEventType_CREATED && event.Price <= 0 {
		return fmt.Errorf("invalid price for CREATED inventory event: %f", event.Price)
	}
	if eventType == pb.InventoryEventType_CREATED && event.Stock <= 0 {
		return fmt.Errorf("invalid stock for CREATED inventory event: %d", event.Stock)
	}

	pbEvent := &pb.InventoryEvent{
		Id:          event.ID.Hex(),
		Name:        event.Name,
		Description: event.Description,
		CategoryId:  event.CategoryID.Hex(),
		Price:       event.Price,
		Stock:       event.Stock,
		CreatedAt:   timestamppb.New(event.CreatedAt),
		UpdatedAt:   timestamppb.New(event.UpdatedAt),
		EventType:   eventType,
	}

	data, err := proto.Marshal(pbEvent)
	if err != nil {
		return fmt.Errorf("proto.Marshal: %w", err)
	}

	log.Printf("Publishing InventoryEvent to %s: %+v, data: %s", p.subject, pbEvent, hex.EncodeToString(data))
	err = p.natsClient.Conn.Publish(p.subject, data)
	if err != nil {
		return fmt.Errorf("p.natsClient.Conn.Publish: %w", err)
	}
	log.Printf("Inventory event pushed to %s: %+v [%s]", p.subject, event, eventType)

	return nil
}
