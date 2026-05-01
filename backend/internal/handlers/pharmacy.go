package handlers

import (
	"encoding/json"
	"net/http"

	"pharmasense/backend/internal/middleware"
	"pharmasense/backend/internal/models"
	"pharmasense/backend/internal/repository"
)

type PharmacyHandler struct {
	repo *repository.Repository
}

func NewPharmacyHandler(repo *repository.Repository) *PharmacyHandler {
	return &PharmacyHandler{repo: repo}
}

func (h *PharmacyHandler) Get(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	p, err := h.repo.GetPharmacy(r.Context(), claims.PharmacyID)
	if err != nil {
		writeError(w, http.StatusNotFound, "pharmacy not found")
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *PharmacyHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	if claims.Role != "admin" {
		writeError(w, http.StatusForbidden, "admin only")
		return
	}

	p, err := h.repo.GetPharmacy(r.Context(), claims.PharmacyID)
	if err != nil {
		writeError(w, http.StatusNotFound, "pharmacy not found")
		return
	}

	var update models.Pharmacy
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}

	p.Name = update.Name
	p.Address = update.Address
	p.City = update.City
	p.Phone = update.Phone
	p.Email = update.Email
	p.Language = update.Language

	if err := h.repo.UpdatePharmacy(r.Context(), p); err != nil {
		writeError(w, http.StatusInternalServerError, "update failed")
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *PharmacyHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	users, err := h.repo.ListUsers(r.Context(), claims.PharmacyID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "list users failed")
		return
	}
	writeJSON(w, http.StatusOK, users)
}
