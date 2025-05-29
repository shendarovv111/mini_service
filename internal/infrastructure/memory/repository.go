package memory

import (
	"context"
	"sync"

	"github.com/anastasiakormilina/quotes-service/internal/domain/quote"
)

type Repository struct {
	mu     sync.RWMutex
	quotes map[int]quote.Quote
	nextID int
}

func NewRepository() *Repository {
	return &Repository{
		quotes: make(map[int]quote.Quote),
		nextID: 1,
	}
}

func (r *Repository) Create(_ context.Context, q quote.Quote) (quote.Quote, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	q.ID = r.nextID
	r.quotes[q.ID] = q
	r.nextID++

	return q, nil
}

func (r *Repository) GetAll(_ context.Context) ([]quote.Quote, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	quotes := make([]quote.Quote, 0, len(r.quotes))
	for _, q := range r.quotes {
		quotes = append(quotes, q)
	}
	return quotes, nil
}

func (r *Repository) GetByAuthor(_ context.Context, author string) ([]quote.Quote, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	quotes := make([]quote.Quote, 0)
	for _, q := range r.quotes {
		if q.Author == author {
			quotes = append(quotes, q)
		}
	}
	return quotes, nil
}

func (r *Repository) Delete(_ context.Context, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.quotes[id]; !exists {
		return quote.ErrQuoteNotFound
	}

	delete(r.quotes, id)
	return nil
}

func (r *Repository) GetByID(_ context.Context, id int) (quote.Quote, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	q, exists := r.quotes[id]
	if !exists {
		return quote.Quote{}, quote.ErrQuoteNotFound
	}
	return q, nil
}
