package handlers

import (
	"net/http"

	"pharmasense/backend/internal/middleware"
	"pharmasense/backend/internal/repository"
	"pharmasense/backend/internal/services"
)

type RiskHandler struct {
	repo   *repository.Repository
	engine *services.RiskEngine
}

func NewRiskHandler(repo *repository.Repository, engine *services.RiskEngine) *RiskHandler {
	return &RiskHandler{repo: repo, engine: engine}
}

func (h *RiskHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	data, err := h.repo.GetDashboard(r.Context(), claims.PharmacyID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "dashboard data failed")
		return
	}
	writeJSON(w, http.StatusOK, data)
}

func (h *RiskHandler) Timeline(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	data, err := h.repo.GetTimeline(r.Context(), claims.PharmacyID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "timeline failed")
		return
	}
	writeJSON(w, http.StatusOK, data)
}

func (h *RiskHandler) TopRisk(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	limit := queryInt(r, "limit", 10)
	data, err := h.repo.GetTopRisk(r.Context(), claims.PharmacyID, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "top risk failed")
		return
	}
	writeJSON(w, http.StatusOK, data)
}

func (h *RiskHandler) Recalculate(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	go func() {
		_ = h.engine.RecalculateAll(r.Context(), claims.PharmacyID)
	}()
	writeJSON(w, http.StatusAccepted, map[string]string{"message": "recalculation started"})
}

func (h *RiskHandler) Assessments(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)

	f := repository.InventoryFilter{
		PharmacyID: claims.PharmacyID,
		RiskLevel:  queryStr(r, "risk_level"),
		Page:       queryInt(r, "page", 1),
		Limit:      queryInt(r, "limit", 50),
	}

	batches, total, err := h.repo.ListInventoryWithRisk(r.Context(), f)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "assessments failed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"data":  batches,
		"total": total,
		"page":  f.Page,
		"limit": f.Limit,
	})
}
