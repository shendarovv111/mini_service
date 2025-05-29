package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/anastasiakormilina/quotes-service/internal/domain/quote"
)

type Handler struct {
	service quote.Service
}

func NewHandler(service quote.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateQuote(w http.ResponseWriter, r *http.Request) {
	var req CreateQuoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	q := req.ToQuote()
	createdQuote, err := h.service.Create(r.Context(), q)
	if err != nil {
		handleError(w, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, FromQuote(createdQuote))
}

func (h *Handler) GetAllQuotes(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author")

	var quotes []quote.Quote
	var err error

	if author != "" {
		quotes, err = h.service.GetByAuthor(r.Context(), author)
	} else {
		quotes, err = h.service.GetAll(r.Context())
	}

	if err != nil {
		handleError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, FromQuotes(quotes))
}

func (h *Handler) GetRandomQuote(w http.ResponseWriter, r *http.Request) {
	q, err := h.service.GetRandom(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, FromQuote(q))
}

func (h *Handler) DeleteQuote(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid quote ID")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, ErrorResponse{Error: message})
}

func handleError(w http.ResponseWriter, err error) {
	code := http.StatusInternalServerError
	message := "Internal server error"

	switch {
	case errors.Is(err, quote.ErrQuoteNotFound):
		code = http.StatusNotFound
		message = "Quote not found"
	case errors.Is(err, quote.ErrNoQuotesAvailable):
		code = http.StatusNotFound
		message = "No quotes available"
	case errors.Is(err, quote.ErrInvalidQuoteID),
		errors.Is(err, quote.ErrEmptyAuthor),
		errors.Is(err, quote.ErrEmptyQuoteText):
		code = http.StatusBadRequest
		message = err.Error()
	}

	respondWithError(w, code, message)
}
