package api

import (
	"encoding/json"
	"errors"
	"markdown-api/internal/documents"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type DocumentHandler struct {
	repository *documents.Repository
	storage    Storage
}

func NewDocumentHandler(
	repository *documents.Repository,
	storage Storage,
) *DocumentHandler {
	return &DocumentHandler{
		repository: repository,
		storage:    storage,
	}
}

func (h *DocumentHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(
		chi.URLParam(r, "id"),
		10,
		64,
	)

	if err != nil {
		http.Error(
			w,
			"invalid id"+err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	doc, err := h.repository.Get(id)

	if errors.Is(err, documents.ErrNotFound) {
		http.Error(
			w,
			"document not found"+err.Error(),
			http.StatusNotFound,
		)
		return
	}

	if err != nil {
		http.Error(
			w,
			"internal error"+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	content, err := h.storage.Load(doc.Filename)

	if err != nil {
		http.Error(
			w,
			"file missing"+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	response := map[string]any{
		"id":       doc.ID,
		"title":    doc.Title,
		"filename": doc.Filename,
		"content":  content,
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(response)
}
