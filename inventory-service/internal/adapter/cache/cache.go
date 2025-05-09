package cache

import (
	"log"
	"sync"
	"time"

	"github.com/mephirious/advanced-programming-2/inventory-service/internal/domain"
)

type ProductCache struct {
	products map[string]domain.Product
	mu       sync.RWMutex
}

func NewProductCache() *ProductCache {
	return &ProductCache{
		products: make(map[string]domain.Product),
	}
}

func (c *ProductCache) Set(product domain.Product) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.products[product.ID.Hex()] = product
}

func (c *ProductCache) Get(id string) (domain.Product, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	p, found := c.products[id]
	return p, found
}

func (c *ProductCache) GetAll() []domain.Product {
	c.mu.RLock()
	defer c.mu.RUnlock()
	all := make([]domain.Product, 0, len(c.products))
	for _, p := range c.products {
		all = append(all, p)
	}
	return all
}

func (c *ProductCache) SetMany(products []domain.Product) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, p := range products {
		c.products[p.ID.Hex()] = p
	}
}

func (c *ProductCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.products = make(map[string]domain.Product)
}

func StartCacheRefresher(cache *ProductCache) {
	go func() {
		for {
			time.Sleep(12 * time.Hour)
			cache.Clear()
			log.Println("Cache cleared after 12 hours")
		}
	}()
}
