package quote

import (
	"context"
	"testing"
	"time"
)

type simpleRepository struct {
	quotes map[int]Quote
	nextID int
}

func newSimpleRepository() *simpleRepository {
	return &simpleRepository{
		quotes: make(map[int]Quote),
		nextID: 1,
	}
}

func (r *simpleRepository) Create(_ context.Context, quote Quote) (Quote, error) {
	quote.ID = r.nextID
	r.quotes[quote.ID] = quote
	r.nextID++
	return quote, nil
}

func (r *simpleRepository) GetAll(_ context.Context) ([]Quote, error) {
	quotes := make([]Quote, 0, len(r.quotes))
	for _, q := range r.quotes {
		quotes = append(quotes, q)
	}
	return quotes, nil
}

func (r *simpleRepository) GetByAuthor(_ context.Context, author string) ([]Quote, error) {
	quotes := make([]Quote, 0)
	for _, q := range r.quotes {
		if q.Author == author {
			quotes = append(quotes, q)
		}
	}
	return quotes, nil
}

func (r *simpleRepository) Delete(_ context.Context, id int) error {
	if _, exists := r.quotes[id]; !exists {
		return ErrQuoteNotFound
	}
	delete(r.quotes, id)
	return nil
}

func (r *simpleRepository) GetByID(_ context.Context, id int) (Quote, error) {
	quote, exists := r.quotes[id]
	if !exists {
		return Quote{}, ErrQuoteNotFound
	}
	return quote, nil
}

func TestService_Create(t *testing.T) {
	repo := newSimpleRepository()
	service := NewService(repo)
	ctx := context.Background()

	tests := []struct {
		name    string
		quote   Quote
		wantErr bool
	}{
		{
			name: "Валидная цитата",
			quote: Quote{
				Author: "Test Author",
				Text:   "Test Quote",
			},
			wantErr: false,
		},
		{
			name: "Пустой автор",
			quote: Quote{
				Author: "",
				Text:   "Test Quote",
			},
			wantErr: true,
		},
		{
			name: "Пустой текст цитаты",
			quote: Quote{
				Author: "Test Author",
				Text:   "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Create(ctx, tt.quote)

			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if result.ID <= 0 {
					t.Error("Ожидался валидный ID, получен нулевой или отрицательный")
				}

				if result.Author != tt.quote.Author {
					t.Errorf("Ожидался автор %s, получен %s", tt.quote.Author, result.Author)
				}

				if result.Text != tt.quote.Text {
					t.Errorf("Ожидался текст %s, получен %s", tt.quote.Text, result.Text)
				}

				if result.Date.IsZero() {
					t.Error("Ожидалась установленная дата, получена нулевая")
				}
			}
		})
	}
}

func TestService_GetRandom(t *testing.T) {
	repo := newSimpleRepository()
	service := NewService(repo)
	ctx := context.Background()

	t.Run("Пустой репозиторий", func(t *testing.T) {
		_, err := service.GetRandom(ctx)
		if err == nil {
			t.Error("Ожидалась ошибка, получен nil")
		}
	})

	quotes := []Quote{
		{Author: "Author 1", Text: "Quote 1", Date: time.Now()},
		{Author: "Author 2", Text: "Quote 2", Date: time.Now()},
		{Author: "Author 3", Text: "Quote 3", Date: time.Now()},
	}

	for _, q := range quotes {
		_, err := repo.Create(ctx, q)
		if err != nil {
			t.Fatalf("Не удалось создать цитату: %v", err)
		}
	}

	t.Run("Получение случайной цитаты", func(t *testing.T) {
		quote, err := service.GetRandom(ctx)
		if err != nil {
			t.Errorf("Неожиданная ошибка: %v", err)
			return
		}

		found := false
		for _, q := range quotes {
			if quote.Author == q.Author && quote.Text == q.Text {
				found = true
				break
			}
		}

		if !found {
			t.Error("Полученная цитата не найдена в репозитории")
		}
	})
}

func TestService_GetByAuthor(t *testing.T) {
	repo := newSimpleRepository()
	service := NewService(repo)
	ctx := context.Background()

	authors := []string{"Author 1", "Author 2", "Author 1"}
	quotes := []string{"Quote 1", "Quote 2", "Quote 3"}

	for i := 0; i < len(authors); i++ {
		repo.Create(ctx, Quote{
			Author: authors[i],
			Text:   quotes[i],
			Date:   time.Now(),
		})
	}

	tests := []struct {
		name      string
		author    string
		wantCount int
		wantErr   bool
	}{
		{
			name:      "Существующий автор с несколькими цитатами",
			author:    "Author 1",
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:      "Существующий автор с одной цитатой",
			author:    "Author 2",
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "Несуществующий автор",
			author:    "Author 3",
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:      "Пустой автор",
			author:    "",
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quotes, err := service.GetByAuthor(ctx, tt.author)

			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetByAuthor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(quotes) != tt.wantCount {
				t.Errorf("Service.GetByAuthor() вернул %d цитат, ожидалось %d", len(quotes), tt.wantCount)
			}
		})
	}
}

func TestService_Delete(t *testing.T) {
	repo := newSimpleRepository()
	service := NewService(repo)
	ctx := context.Background()

	quote, _ := repo.Create(ctx, Quote{
		Author: "Test Author",
		Text:   "Test Quote",
		Date:   time.Now(),
	})

	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{
			name:    "Существующий ID",
			id:      quote.ID,
			wantErr: false,
		},
		{
			name:    "Несуществующий ID",
			id:      999,
			wantErr: true,
		},
		{
			name:    "Некорректный ID",
			id:      0,
			wantErr: true,
		},
		{
			name:    "Отрицательный ID",
			id:      -1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.Delete(ctx, tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				_, err := repo.GetByID(ctx, tt.id)
				if err == nil {
					t.Error("Цитата не была удалена из репозитория")
				}
			}
		})
	}
}
