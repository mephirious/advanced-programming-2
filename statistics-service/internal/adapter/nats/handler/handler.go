package handler

import (
	"context"
	"encoding/json"
	"log"

	"github.com/mephirious/advanced-programming-2/statistics-service/internal/domain/dto"
	"github.com/mephirious/advanced-programming-2/statistics-service/internal/usecase"
	"github.com/nats-io/nats.go"
)

type NATSHandler struct {
	uc usecase.StatsUseCase
	nc *nats.Conn
}

func NewNATSHandler(nc *nats.Conn, uc usecase.StatsUseCase) *NATSHandler {
	return &NATSHandler{
		uc: uc,
		nc: nc,
	}
}

func (h *NATSHandler) Start() {
	h.nc.Subscribe("order.*", h.handleOrderEvent)
	h.nc.Subscribe("inventory.*", h.handleInventoryEvent)
}

func (h *NATSHandler) handleOrderEvent(msg *nats.Msg) {
	ctx := context.Background()

	var eventDTO dto.OrderEventDTO
	if err := json.Unmarshal(msg.Data, &eventDTO); err != nil {
		log.Printf("Error unmarshaling order event: %v", err)
		return
	}

	event := dto.ToDomainOrderEvent(eventDTO)
	if err := h.uc.HandleOrderEvent(ctx, &event); err != nil {
		log.Printf("Error processing order event: %v", err)
	}
}

func (h *NATSHandler) handleInventoryEvent(msg *nats.Msg) {
	ctx := context.Background()

	var eventDTO dto.InventoryEventDTO
	if err := json.Unmarshal(msg.Data, &eventDTO); err != nil {
		log.Printf("Error unmarshaling inventory event: %v", err)
		return
	}

	event := dto.ToDomainInventoryEvent(eventDTO)
	if err := h.uc.HandleInventoryEvent(ctx, &event); err != nil {
		log.Printf("Error processing inventory event: %v", err)
	}
}
