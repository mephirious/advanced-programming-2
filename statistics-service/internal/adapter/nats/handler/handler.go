package handler

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/mephirious/advanced-programming-2/statistics-service/internal/domain"
	"github.com/mephirious/advanced-programming-2/statistics-service/internal/usecase"
	pb "github.com/mephirious/advanced-programming-2/statistics-service/proto"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type NATSHandler struct {
	statsUseCase usecase.StatsUseCase
	natsConn     *nats.Conn
}

func NewNATSHandler(statsUC usecase.StatsUseCase, nc *nats.Conn) *NATSHandler {
	return &NATSHandler{
		statsUseCase: statsUC,
		natsConn:     nc,
	}
}

func (h *NATSHandler) Start() error {
	_, err := h.natsConn.Subscribe("order.events", func(m *nats.Msg) {
		log.Printf("Received on order.events: %s", hex.EncodeToString(m.Data))
		var orderEvent pb.OrderEvent
		if err := proto.Unmarshal(m.Data, &orderEvent); err != nil {
			log.Printf("Failed to unmarshal order event: %v", err)
			return
		}
		log.Printf("Successfully unmarshaled OrderEvent: %+v", &orderEvent)
		h.processOrderEvent(context.Background(), &orderEvent)
	})
	if err != nil {
		log.Printf("Error subscribing to subject order.events: %v", err)
		return fmt.Errorf("failed to subscribe to subject order.events: %w", err)
	}
	log.Printf("Subscribed to NATS subject: order.events")

	_, err = h.natsConn.Subscribe("inventory.events", func(m *nats.Msg) {
		log.Printf("Received on inventory.events: %s", hex.EncodeToString(m.Data))
		var inventoryEvent pb.InventoryEvent
		if err := proto.Unmarshal(m.Data, &inventoryEvent); err != nil {
			log.Printf("Failed to unmarshal inventory event: %v", err)
			var orderEvent pb.OrderEvent
			if err := proto.Unmarshal(m.Data, &orderEvent); err == nil {
				log.Printf("Message is actually an OrderEvent: %+v", orderEvent)
			}
			return
		}
		log.Printf("Successfully unmarshaled InventoryEvent: id=%s, name=%s, description=%s, category_id=%s, price=%f, quantity=%d, event_type=%s",
			inventoryEvent.Id, inventoryEvent.Name, inventoryEvent.Description, inventoryEvent.CategoryId, inventoryEvent.Price, inventoryEvent.Quantity, inventoryEvent.EventType)
		h.processInventoryEvent(context.Background(), &inventoryEvent)
	})
	if err != nil {
		log.Printf("Error subscribing to subject inventory.events: %v", err)
		return fmt.Errorf("failed to subscribe to subject inventory.events: %w", err)
	}
	log.Printf("Subscribed to NATS subject: inventory.events")

	return nil
}

func (h *NATSHandler) processOrderEvent(ctx context.Context, pbEvent *pb.OrderEvent) {
	var createdAt, updatedAt time.Time
	if pbEvent.CreatedAt != nil {
		createdAt = pbEvent.CreatedAt.AsTime()
	}
	if pbEvent.UpdatedAt != nil {
		updatedAt = pbEvent.UpdatedAt.AsTime()
	}

	domainEvent := &domain.OrderEvent{
		ID:        pbEvent.Id,
		UserID:    pbEvent.UserId,
		Total:     pbEvent.Total,
		Status:    pbEvent.Status.String(),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		EventType: mapEventTypeToString(pbEvent.EventType),
	}

	for _, item := range pbEvent.Items {
		domainEvent.Items = append(domainEvent.Items, domain.OrderItem{
			ProductID: item.ProductId,
			Quantity:  int(item.Quantity),
			Price:     item.Price,
		})
	}

	h.handleOrderEvent(ctx, domainEvent)
}

func (h *NATSHandler) processInventoryEvent(ctx context.Context, pbEvent *pb.InventoryEvent) {
	var createdAt, updatedAt time.Time
	if pbEvent.CreatedAt != nil {
		createdAt = pbEvent.CreatedAt.AsTime()
	}
	if pbEvent.UpdatedAt != nil {
		updatedAt = pbEvent.UpdatedAt.AsTime()
	}

	domainEvent := &domain.InventoryEvent{
		ID:          pbEvent.Id,
		Name:        pbEvent.Name,
		Description: pbEvent.Description,
		CategoryID:  pbEvent.CategoryId,
		Price:       float64(pbEvent.Price),
		Quantity:    int(pbEvent.Quantity),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		EventType:   mapEventTypeToString(pbEvent.EventType),
	}

	h.handleInventoryEvent(ctx, domainEvent)
}

func (h *NATSHandler) handleOrderEvent(ctx context.Context, event *domain.OrderEvent) {
	if event.EventType == "UNKNOWN" {
		log.Printf("Unknown event type for order event %s", event.ID)
		return
	}

	if err := h.statsUseCase.HandleOrderEvent(ctx, event); err != nil {
		log.Printf("Error handling order event: %v", err)
	}
}

func (h *NATSHandler) handleInventoryEvent(ctx context.Context, event *domain.InventoryEvent) {
	if event.EventType == "UNKNOWN" {
		log.Printf("Unknown event type for inventory event %s", event.ID)
		return
	}

	if event.EventType == "CREATED" {
		if event.Price <= 0 {
			log.Printf("Invalid price for CREATED inventory event %s: %f", event.ID, event.Price)
			return
		}
		if event.Quantity <= 0 {
			log.Printf("Invalid quantity for CREATED inventory event %s: %d", event.ID, event.Quantity)
			return
		}
		if event.CreatedAt.IsZero() {
			log.Printf("CreatedAt is zero for CREATED inventory event %s", event.ID)
			return
		}
	}

	if err := h.statsUseCase.HandleInventoryEvent(ctx, event); err != nil {
		log.Printf("Error handling inventory event: %v", err)
	}
}

func mapEventTypeToString(eventType pb.OrderEventType) string {
	switch eventType {
	case pb.OrderEventType_CREATED:
		return "CREATED"
	case pb.OrderEventType_UPDATED:
		return "UPDATED"
	case pb.OrderEventType_CANCELLED:
		return "CANCELLED"
	case pb.OrderEventType_DELETED:
		return "DELETED"
	default:
		return "UNKNOWN"
	}
}
