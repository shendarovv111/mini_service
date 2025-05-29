package quote

import (
	"errors"
)

var (
	ErrQuoteNotFound     = errors.New("quote not found")
	ErrInvalidQuoteID    = errors.New("invalid quote id")
	ErrEmptyAuthor       = errors.New("author cannot be empty")
	ErrEmptyQuoteText    = errors.New("quote text cannot be empty")
	ErrNoQuotesAvailable = errors.New("no quotes available")
)
