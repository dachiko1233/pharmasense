package repository

import (
	"context"

	"github.com/google/uuid"
)

// AvgDailySales returns the rolling 90-day average daily sales for a product in a pharmacy.
func (r *Repository) AvgDailySales(ctx context.Context, pharmacyID, productID uuid.UUID) (float64, error) {
	var avg float64
	err := r.pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(quantity)::float / 90.0, 0)
		FROM sales
		WHERE pharmacy_id = $1 AND product_id = $2
		  AND sale_date >= CURRENT_DATE - INTERVAL '90 days'`,
		pharmacyID, productID).Scan(&avg)
	return avg, err
}

// AvgDailySalesMap returns avg_daily_sales for all products in a pharmacy in one query.
func (r *Repository) AvgDailySalesMap(ctx context.Context, pharmacyID uuid.UUID) (map[uuid.UUID]float64, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT product_id, COALESCE(SUM(quantity)::float / 90.0, 0) as avg_daily
		FROM sales
		WHERE pharmacy_id = $1
		  AND sale_date >= CURRENT_DATE - INTERVAL '90 days'
		GROUP BY product_id`, pharmacyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[uuid.UUID]float64)
	for rows.Next() {
		var pid uuid.UUID
		var avg float64
		if err := rows.Scan(&pid, &avg); err != nil {
			return nil, err
		}
		result[pid] = avg
	}
	return result, nil
}
