package handlers

import (
	"net/http"

	"pharmasense/backend/internal/middleware"
	"pharmasense/backend/internal/repository"
)

type ReportsHandler struct {
	repo *repository.Repository
}

func NewReportsHandler(repo *repository.Repository) *ReportsHandler {
	return &ReportsHandler{repo: repo}
}

func (h *ReportsHandler) Savings(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	data, err := h.repo.GetSavingsReport(r.Context(), claims.PharmacyID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "savings report failed")
		return
	}
	writeJSON(w, http.StatusOK, data)
}

func (h *ReportsHandler) Categories(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	data, err := h.repo.GetCategoryReport(r.Context(), claims.PharmacyID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "category report failed")
		return
	}
	writeJSON(w, http.StatusOK, data)
}

func (h *ReportsHandler) Waste(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	data, err := h.repo.GetSavingsReport(r.Context(), claims.PharmacyID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "waste report failed")
		return
	}
	writeJSON(w, http.StatusOK, data)
}
