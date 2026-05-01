package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"pharmasense/backend/internal/middleware"
	"pharmasense/backend/internal/models"
	"pharmasense/backend/internal/repository"
)

type AlertsHandler struct {
	repo *repository.Repository
}

func NewAlertsHandler(repo *repository.Repository) *AlertsHandler {
	return &AlertsHandler{repo: repo}
}

func (h *AlertsHandler) List(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)

	riskLevel := queryStr(r, "risk_level")
	if riskLevel == "" {
		riskLevel = "" // all CRITICAL and HIGH
	}

	f := repository.InventoryFilter{
		PharmacyID: claims.PharmacyID,
		RiskLevel:  riskLevel,
		Page:       queryInt(r, "page", 1),
		Limit:      queryInt(r, "limit", 50),
	}

	// Default to critical/high alerts
	if f.RiskLevel == "" {
		// fetch critical first, then high — combine
		critF := f
		critF.RiskLevel = "CRITICAL"
		critF.Limit = 100
		crits, _, _ := h.repo.ListInventoryWithRisk(r.Context(), critF)

		highF := f
		highF.RiskLevel = "HIGH"
		highF.Limit = 100
		highs, _, _ := h.repo.ListInventoryWithRisk(r.Context(), highF)

		all := append(crits, highs...)
		writeJSON(w, http.StatusOK, map[string]any{
			"data":  all,
			"total": len(all),
		})
		return
	}

	batches, total, err := h.repo.ListInventoryWithRisk(r.Context(), f)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "list alerts failed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"data":  batches,
		"total": total,
	})
}

func (h *AlertsHandler) TakeAction(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)

	batchID, err := uuid.Parse(chi.URLParam(r, "batch_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid batch_id")
		return
	}

	var body struct {
		ActionType      string `json:"action_type"`
		DiscountPercent *int   `json:"discount_percent"`
		Notes           string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}

	validActions := map[string]bool{"DISCOUNT": true, "TRANSFER": true, "RETURN": true, "DISMISS": true}
	if !validActions[body.ActionType] {
		writeError(w, http.StatusBadRequest, "invalid action_type")
		return
	}

	action := &models.AlertAction{
		BatchID:         batchID,
		PharmacyID:      claims.PharmacyID,
		UserID:          claims.UserID,
		ActionType:      body.ActionType,
		DiscountPercent: body.DiscountPercent,
		Notes:           body.Notes,
	}

	if err := h.repo.CreateAlertAction(r.Context(), action); err != nil {
		writeError(w, http.StatusInternalServerError, "action failed")
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"message": "action recorded"})
}

func (h *AlertsHandler) History(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	limit := queryInt(r, "limit", 50)

	actions, err := h.repo.ListAlertActions(r.Context(), claims.PharmacyID, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "history failed")
		return
	}
	writeJSON(w, http.StatusOK, actions)
}
