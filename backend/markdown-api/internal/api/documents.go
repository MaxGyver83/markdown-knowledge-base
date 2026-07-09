package api

import (
	"encoding/json"
	"errors"
	"fmt"
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
			err.Error(),
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

	response := DocumentResponse{
		DocumentMetadata: DocumentMetadata{
			ID:        doc.ID,
			Title:     doc.Title,
			CreatedAt: doc.CreatedAt,
			UpdatedAt: doc.UpdatedAt,
		},
		Content: content,
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(response)
}

func (h *DocumentHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req CreateDocumentRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(
			w,
			"invalid json"+err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	doc := documents.Document{
		Title: req.Title,
	}

	err = h.repository.Create(&doc)

	if err != nil {
		http.Error(
			w,
			"could not create document: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	doc.Filename = fmt.Sprintf(
		"%d.md",
		doc.ID,
	)

	err = h.storage.Save(
		doc.Filename,
		req.Content,
	)

	if err != nil {
		http.Error(
			w,
			"could not save file"+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	err = h.repository.UpdateFilename(
		doc.ID,
		doc.Filename,
	)

	if err != nil {
		http.Error(
			w,
			"could not update document"+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(doc)
}

func (h *DocumentHandler) List(
	w http.ResponseWriter,
	r *http.Request,
) {
	docs, err := h.repository.List()

	if err != nil {
		http.Error(
			w,
			"could not list documents: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	response := make([]DocumentMetadata, 0, len(docs))

	for _, doc := range docs {
		response = append(response, DocumentMetadata{
			ID:        doc.ID,
			Title:     doc.Title,
			CreatedAt: doc.CreatedAt,
			UpdatedAt: doc.UpdatedAt,
		})
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(
			w,
			"could not encode response",
			http.StatusInternalServerError,
		)
	}
}
