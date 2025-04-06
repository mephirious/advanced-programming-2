package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type ProductResponse struct {
	ID    string  `json:"id"`
	Price float64 `json:"price"`
}

type ProductGateway interface {
	GetProductPrices(ctx context.Context, productIDs []string) (map[string]float64, error)
}

type HTTPGateway struct {
	baseURL string
	client  *http.Client
}

func NewHTTPGateway(baseURL string) *HTTPGateway {
	return &HTTPGateway{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (g *HTTPGateway) GetProductPrices(ctx context.Context, productIDs []string) (map[string]float64, error) {
	if len(productIDs) == 0 {
		return nil, fmt.Errorf("no product IDs provided")
	}

	idsParam := strings.Join(productIDs, ",")
	url := fmt.Sprintf("%s/api/v1/products?ids=%s", g.baseURL, idsParam)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var products []ProductResponse
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	prices := make(map[string]float64)
	for _, p := range products {
		prices[p.ID] = p.Price
	}

	for _, id := range productIDs {
		if _, exists := prices[id]; !exists {
			return nil, fmt.Errorf("price not found for product %s", id)
		}
	}

	return prices, nil
}
