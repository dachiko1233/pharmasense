package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"pharmasense/backend/internal/models"
	"pharmasense/backend/internal/repository"
)

type RiskEngine struct {
	repo *repository.Repository
}

func NewRiskEngine(repo *repository.Repository) *RiskEngine {
	return &RiskEngine{repo: repo}
}

// Calculate computes the risk level for a single batch given avg daily sales.
func Calculate(daysUntilExpiry int, currentQty int, avgDailySales float64, purchasePrice float64) models.RiskAssessment {
	expectedSales := int(avgDailySales * float64(daysUntilExpiry))
	surplus := currentQty - expectedSales
	if surplus < 0 {
		surplus = 0
	}

	var level models.RiskLevel
	var discountPct int
	var estimatedLoss float64

	switch {
	case daysUntilExpiry <= 30 && surplus > 0:
		level = models.RiskCritical
		discountPct = 40
	case daysUntilExpiry <= 90 && expectedSales > 0 && float64(surplus) > float64(expectedSales)*0.5:
		level = models.RiskHigh
		discountPct = 20
	case daysUntilExpiry <= 90 && surplus > 0 && expectedSales == 0:
		// Zero sales, expiring soon — treat as high
		level = models.RiskHigh
		discountPct = 20
	case daysUntilExpiry <= 180 && expectedSales > 0 && float64(surplus) > float64(expectedSales)*0.3:
		level = models.RiskMedium
		discountPct = 10
	case daysUntilExpiry <= 0:
		level = models.RiskCritical
		discountPct = 50
		surplus = currentQty
	default:
		level = models.RiskLow
	}

	if level == models.RiskCritical || level == models.RiskHigh {
		estimatedLoss = float64(surplus) * purchasePrice
	}
	if level == models.RiskMedium {
		estimatedLoss = float64(surplus) * purchasePrice * 0.3
	}

	return models.RiskAssessment{
		RiskLevel:            level,
		DaysUntilExpiry:      daysUntilExpiry,
		AvgDailySales:        avgDailySales,
		ExpectedSales:        expectedSales,
		EstimatedSurplus:     surplus,
		EstimatedLoss:        estimatedLoss,
		SuggestedDiscountPct: discountPct,
	}
}

// RecalculateAll recomputes risk assessments for every batch in a pharmacy.
func (e *RiskEngine) RecalculateAll(ctx context.Context, pharmacyID uuid.UUID) error {
	slog.Info("recalculating risk", "pharmacy_id", pharmacyID)

	salesMap, err := e.repo.AvgDailySalesMap(ctx, pharmacyID)
	if err != nil {
		return err
	}

	batches, err := e.repo.ListBatchIDsForPharmacy(ctx, pharmacyID)
	if err != nil {
		return err
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)

	for _, b := range batches {
		expiryDate, err := time.Parse("2006-01-02", b.ExpiryDate)
		if err != nil {
			continue
		}
		days := int(expiryDate.Sub(today).Hours() / 24)
		avg := salesMap[b.ProductID]

		ra := Calculate(days, b.CurrentQty, avg, b.PurchasePrice)
		ra.BatchID = b.BatchID
		ra.PharmacyID = b.PharmacyID

		if err := e.repo.UpsertRiskAssessment(ctx, &ra); err != nil {
			slog.Error("upsert risk", "batch_id", b.BatchID, "error", err)
		}
	}

	slog.Info("risk recalculation complete", "batches", len(batches))
	return nil
}
