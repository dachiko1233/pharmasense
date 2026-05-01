package repository

import (
	"context"
	"fmt"

	"pharmasense/backend/internal/models"

	"github.com/google/uuid"
)

func (r *Repository) UpsertRiskAssessment(ctx context.Context, ra *models.RiskAssessment) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO risk_assessments (
			id, batch_id, pharmacy_id, risk_level, days_until_expiry,
			avg_daily_sales, expected_sales, estimated_surplus,
			estimated_loss, suggested_discount_percent, calculated_at
		) VALUES (
			gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, $9, NOW()
		)
		ON CONFLICT (batch_id) DO UPDATE SET
			risk_level = EXCLUDED.risk_level,
			days_until_expiry = EXCLUDED.days_until_expiry,
			avg_daily_sales = EXCLUDED.avg_daily_sales,
			expected_sales = EXCLUDED.expected_sales,
			estimated_surplus = EXCLUDED.estimated_surplus,
			estimated_loss = EXCLUDED.estimated_loss,
			suggested_discount_percent = EXCLUDED.suggested_discount_percent,
			calculated_at = NOW()`,
		ra.BatchID, ra.PharmacyID, string(ra.RiskLevel), ra.DaysUntilExpiry,
		ra.AvgDailySales, ra.ExpectedSales, ra.EstimatedSurplus,
		ra.EstimatedLoss, ra.SuggestedDiscountPct)
	return err
}

func (r *Repository) GetDashboard(ctx context.Context, pharmacyID uuid.UUID) (*models.DashboardData, error) {
	d := &models.DashboardData{}

	// Risk counts
	rows, err := r.pool.Query(ctx, `
		SELECT risk_level, COUNT(*), COALESCE(SUM(estimated_loss),0)
		FROM risk_assessments
		WHERE pharmacy_id = $1
		GROUP BY risk_level`, pharmacyID)
	if err != nil {
		return nil, fmt.Errorf("dashboard counts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var level string
		var count int
		var loss float64
		if err := rows.Scan(&level, &count, &loss); err != nil {
			return nil, err
		}
		switch level {
		case "CRITICAL":
			d.CriticalCount = count
			d.TotalEstimatedLoss += loss
		case "HIGH":
			d.HighCount = count
			d.TotalEstimatedLoss += loss
		case "MEDIUM":
			d.MediumCount = count
		case "LOW":
			d.LowCount = count
		}
	}

	d.TotalPotentialSavings = d.TotalEstimatedLoss * 0.6

	// Total inventory value
	err = r.pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(current_quantity * selling_price), 0)
		FROM inventory_batches
		WHERE pharmacy_id = $1`, pharmacyID).Scan(&d.TotalInventoryValue)
	if err != nil {
		return nil, err
	}

	// Top 5 critical alerts
	alertRows, err := r.pool.Query(ctx, `
		SELECT ib.id, ib.pharmacy_id, ib.product_id, p.name, COALESCE(p.name_el,''), COALESCE(p.category,''),
		       COALESCE(p.manufacturer,''), COALESCE(p.barcode,''), COALESCE(ib.batch_number,''),
		       ib.expiry_date, ib.initial_quantity, ib.current_quantity,
		       ib.purchase_price, ib.selling_price, COALESCE(ib.supplier,''), ib.received_date,
		       ra.risk_level, ra.days_until_expiry, COALESCE(ra.avg_daily_sales,0),
		       COALESCE(ra.expected_sales,0), COALESCE(ra.estimated_surplus,0),
		       COALESCE(ra.estimated_loss,0), COALESCE(ra.suggested_discount_percent,0)
		FROM risk_assessments ra
		JOIN inventory_batches ib ON ib.id = ra.batch_id
		JOIN products p ON p.id = ib.product_id
		WHERE ra.pharmacy_id = $1 AND ra.risk_level = 'CRITICAL'
		ORDER BY ra.estimated_loss DESC
		LIMIT 5`, pharmacyID)
	if err != nil {
		return nil, fmt.Errorf("critical alerts: %w", err)
	}
	defer alertRows.Close()

	for alertRows.Next() {
		var b models.InventoryBatchWithRisk
		if err := alertRows.Scan(
			&b.ID, &b.PharmacyID, &b.ProductID, &b.ProductName, &b.ProductNameEl, &b.Category,
			&b.Manufacturer, &b.Barcode, &b.BatchNumber,
			&b.ExpiryDate, &b.InitialQuantity, &b.CurrentQuantity,
			&b.PurchasePrice, &b.SellingPrice, &b.Supplier, &b.ReceivedDate,
			&b.RiskLevel, &b.DaysUntilExpiry, &b.AvgDailySales, &b.ExpectedSales,
			&b.EstimatedSurplus, &b.EstimatedLoss, &b.SuggestedDiscountPct,
		); err != nil {
			return nil, err
		}
		d.CriticalAlerts = append(d.CriticalAlerts, b)
	}

	return d, nil
}

func (r *Repository) GetTimeline(ctx context.Context, pharmacyID uuid.UUID) ([]models.TimelineMonth, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT
			TO_CHAR(DATE_TRUNC('month', ib.expiry_date), 'YYYY-MM') AS month,
			COALESCE(ra.risk_level, 'LOW') AS risk_level,
			COUNT(*) AS cnt
		FROM inventory_batches ib
		LEFT JOIN risk_assessments ra ON ra.batch_id = ib.id
		WHERE ib.pharmacy_id = $1
		  AND ib.expiry_date BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '12 months'
		GROUP BY 1, 2
		ORDER BY 1, 2`, pharmacyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	monthMap := make(map[string]*models.TimelineMonth)
	var order []string

	for rows.Next() {
		var month, level string
		var cnt int
		if err := rows.Scan(&month, &level, &cnt); err != nil {
			return nil, err
		}
		if _, ok := monthMap[month]; !ok {
			monthMap[month] = &models.TimelineMonth{Month: month}
			order = append(order, month)
		}
		tm := monthMap[month]
		switch level {
		case "CRITICAL":
			tm.Critical += cnt
		case "HIGH":
			tm.High += cnt
		case "MEDIUM":
			tm.Medium += cnt
		default:
			tm.Low += cnt
		}
	}

	result := make([]models.TimelineMonth, 0, len(order))
	for _, m := range order {
		result = append(result, *monthMap[m])
	}
	return result, nil
}

func (r *Repository) GetTopRisk(ctx context.Context, pharmacyID uuid.UUID, limit int) ([]models.TopRiskItem, error) {
	if limit == 0 {
		limit = 10
	}
	rows, err := r.pool.Query(ctx, `
		SELECT p.name, COALESCE(ib.batch_number,''), COALESCE(ra.estimated_loss,0),
		       ra.risk_level, COALESCE(p.category,'')
		FROM risk_assessments ra
		JOIN inventory_batches ib ON ib.id = ra.batch_id
		JOIN products p ON p.id = ib.product_id
		WHERE ra.pharmacy_id = $1 AND ra.risk_level IN ('CRITICAL','HIGH')
		ORDER BY ra.estimated_loss DESC
		LIMIT $2`, pharmacyID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.TopRiskItem
	for rows.Next() {
		var item models.TopRiskItem
		var level string
		if err := rows.Scan(&item.ProductName, &item.BatchNumber, &item.EstimatedLoss, &level, &item.Category); err != nil {
			return nil, err
		}
		item.RiskLevel = models.RiskLevel(level)
		items = append(items, item)
	}
	return items, nil
}

func (r *Repository) ListBatchIDsForPharmacy(ctx context.Context, pharmacyID uuid.UUID) ([]struct {
	BatchID       uuid.UUID
	ProductID     uuid.UUID
	PharmacyID    uuid.UUID
	ExpiryDate    string
	CurrentQty    int
	PurchasePrice float64
}, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, product_id, pharmacy_id, expiry_date::text, current_quantity, purchase_price
		FROM inventory_batches
		WHERE pharmacy_id = $1`, pharmacyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type row []struct {
		BatchID       uuid.UUID
		ProductID     uuid.UUID
		PharmacyID    uuid.UUID
		ExpiryDate    string
		CurrentQty    int
		PurchasePrice float64
	}

	var result []row

	for rows.Next() {
		var r row
		if err := rows.Scan(&r.BatchID, &r.ProductID, &r.PharmacyID, &r.ExpiryDate, &r.CurrentQty, &r.PurchasePrice); err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return result, nil
}
