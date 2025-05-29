package api

import (
	"github.com/anastasiakormilina/quotes-service/internal/domain/quote"
)

type CreateQuoteRequest struct {
	Author string `json:"author"`
	Quote  string `json:"quote"`
}

type QuoteResponse struct {
	ID     int    `json:"id"`
	Author string `json:"author"`
	Quote  string `json:"quote"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (r *CreateQuoteRequest) ToQuote() quote.Quote {
	return quote.Quote{
		Author: r.Author,
		Text:   r.Quote,
	}
}

func FromQuote(q quote.Quote) QuoteResponse {
	return QuoteResponse{
		ID:     q.ID,
		Author: q.Author,
		Quote:  q.Text,
	}
}

func FromQuotes(quotes []quote.Quote) []QuoteResponse {
	result := make([]QuoteResponse, len(quotes))
	for i, q := range quotes {
		result[i] = FromQuote(q)
	}
	return result
}
