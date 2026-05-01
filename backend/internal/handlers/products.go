package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"pharmasense/backend/internal/repository"
)

type ProductsHandler struct {
	repo *repository.Repository
}

func NewProductsHandler(repo *repository.Repository) *ProductsHandler {
	return &ProductsHandler{repo: repo}
}

func (h *ProductsHandler) List(w http.ResponseWriter, r *http.Request) {
	f := repository.ProductFilter{
		Search:   queryStr(r, "search"),
		Category: queryStr(r, "category"),
		Page:     queryInt(r, "page", 1),
		Limit:    queryInt(r, "limit", 50),
	}

	products, total, err := h.repo.ListProducts(r.Context(), f)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "list products failed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"data":  products,
		"total": total,
		"page":  f.Page,
		"limit": f.Limit,
	})
}

func (h *ProductsHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	p, err := h.repo.GetProduct(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "product not found")
		return
	}
	writeJSON(w, http.StatusOK, p)
}
