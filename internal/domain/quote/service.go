package quote

import (
	"context"
	"errors"
	"math/rand"
	"time"
)

type Quote struct {
	ID     int       `json:"id"`
	Author string    `json:"author"`
	Text   string    `json:"quote"`
	Date   time.Time `json:"date"`
}

type Repository interface {
	Create(ctx context.Context, quote Quote) (Quote, error)
	GetAll(ctx context.Context) ([]Quote, error)
	GetByAuthor(ctx context.Context, author string) ([]Quote, error)
	Delete(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (Quote, error)
}

type Service interface {
	Create(ctx context.Context, quote Quote) (Quote, error)
	GetAll(ctx context.Context) ([]Quote, error)
	GetRandom(ctx context.Context) (Quote, error)
	GetByAuthor(ctx context.Context, author string) ([]Quote, error)
	Delete(ctx context.Context, id int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, quote Quote) (Quote, error) {
	if quote.Author == "" || quote.Text == "" {
		return Quote{}, errors.New("author and quote text are required")
	}

	quote.Date = time.Now()
	return s.repo.Create(ctx, quote)
}

func (s *service) GetAll(ctx context.Context) ([]Quote, error) {
	return s.repo.GetAll(ctx)
}

func (s *service) GetRandom(ctx context.Context) (Quote, error) {
	quotes, err := s.repo.GetAll(ctx)
	if err != nil {
		return Quote{}, err
	}

	if len(quotes) == 0 {
		return Quote{}, errors.New("no quotes available")
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return quotes[r.Intn(len(quotes))], nil
}

func (s *service) GetByAuthor(ctx context.Context, author string) ([]Quote, error) {
	if author == "" {
		return nil, errors.New("author cannot be empty")
	}

	return s.repo.GetByAuthor(ctx, author)
}

func (s *service) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid quote ID")
	}

	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}
