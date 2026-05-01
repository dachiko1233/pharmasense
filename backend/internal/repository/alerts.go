package repository

import (
	"context"

	"github.com/google/uuid"
	"pharmasense/backend/internal/models"
)

func (r *Repository) CreateAlertAction(ctx context.Context, a *models.AlertAction) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO alert_actions (id, batch_id, pharmacy_id, user_id, action_type, discount_percent, notes)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6)`,
		a.BatchID, a.PharmacyID, a.UserID, a.ActionType, a.DiscountPercent, a.Notes)
	return err
}

func (r *Repository) ListAlertActions(ctx context.Context, pharmacyID uuid.UUID, limit int) ([]models.AlertActionWithDetails, error) {
	if limit == 0 {
		limit = 50
	}
	rows, err := r.pool.Query(ctx, `
		SELECT aa.id, aa.batch_id, aa.pharmacy_id, aa.user_id, aa.action_type,
		       aa.discount_percent, COALESCE(aa.notes,''), aa.created_at,
		       u.full_name, p.name, COALESCE(ib.batch_number,'')
		FROM alert_actions aa
		JOIN users u ON u.id = aa.user_id
		JOIN inventory_batches ib ON ib.id = aa.batch_id
		JOIN products p ON p.id = ib.product_id
		WHERE aa.pharmacy_id = $1
		ORDER BY aa.created_at DESC
		LIMIT $2`, pharmacyID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actions []models.AlertActionWithDetails
	for rows.Next() {
		var a models.AlertActionWithDetails
		if err := rows.Scan(
			&a.ID, &a.BatchID, &a.PharmacyID, &a.UserID, &a.ActionType,
			&a.DiscountPercent, &a.Notes, &a.CreatedAt,
			&a.UserName, &a.ProductName, &a.BatchNumber,
		); err != nil {
			return nil, err
		}
		actions = append(actions, a)
	}
	return actions, nil
}

func (r *Repository) GetSavingsReport(ctx context.Context, pharmacyID uuid.UUID) ([]models.SavingsReport, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT TO_CHAR(DATE_TRUNC('month', aa.created_at), 'YYYY-MM') as month,
		       COUNT(*) as actions,
		       COALESCE(SUM(
		           CASE WHEN aa.action_type IN ('DISCOUNT','TRANSFER','RETURN')
		           THEN COALESCE(ra.estimated_loss * 0.6, 0) ELSE 0 END
		       ), 0) as estimated_saved
		FROM alert_actions aa
		LEFT JOIN risk_assessments ra ON ra.batch_id = aa.batch_id
		WHERE aa.pharmacy_id = $1
		GROUP BY 1
		ORDER BY 1`, pharmacyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []models.SavingsReport
	for rows.Next() {
		var s models.SavingsReport
		if err := rows.Scan(&s.Month, &s.ActionsCount, &s.EstimatedSaved); err != nil {
			return nil, err
		}
		reports = append(reports, s)
	}
	return reports, nil
}

func (r *Repository) GetCategoryReport(ctx context.Context, pharmacyID uuid.UUID) ([]models.CategoryReport, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT COALESCE(p.category,'Other') as category,
		       COUNT(DISTINCT ib.id) as total_batches,
		       COUNT(DISTINCT CASE WHEN ra.risk_level IN ('CRITICAL','HIGH') THEN ib.id END) as at_risk,
		       COALESCE(SUM(CASE WHEN ra.risk_level IN ('CRITICAL','HIGH') THEN ra.estimated_loss ELSE 0 END),0) as loss
		FROM inventory_batches ib
		JOIN products p ON p.id = ib.product_id
		LEFT JOIN risk_assessments ra ON ra.batch_id = ib.id
		WHERE ib.pharmacy_id = $1
		GROUP BY 1
		ORDER BY 4 DESC`, pharmacyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []models.CategoryReport
	for rows.Next() {
		var c models.CategoryReport
		if err := rows.Scan(&c.Category, &c.TotalBatches, &c.AtRisk, &c.EstimatedLoss); err != nil {
			return nil, err
		}
		reports = append(reports, c)
	}
	return reports, nil
}
