package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"pharmasense/backend/internal/models"
)

type InventoryFilter struct {
	PharmacyID  uuid.UUID
	RiskLevel   string
	Category    string
	Supplier    string
	Search      string
	ExpiryFrom  *time.Time
	ExpiryTo    *time.Time
	Page        int
	Limit       int
	SortBy      string
	SortDir     string
}

func (r *Repository) ListInventoryWithRisk(ctx context.Context, f InventoryFilter) ([]models.InventoryBatchWithRisk, int, error) {
	if f.Limit == 0 {
		f.Limit = 50
	}
	if f.Page == 0 {
		f.Page = 1
	}
	if f.SortBy == "" {
		f.SortBy = "ra.risk_level"
	}

	var conds []string
	var args []any
	args = append(args, f.PharmacyID)
	i := 2

	conds = append(conds, "ib.pharmacy_id = $1")

	if f.RiskLevel != "" {
		conds = append(conds, fmt.Sprintf("ra.risk_level = $%d", i))
		args = append(args, f.RiskLevel)
		i++
	}
	if f.Category != "" {
		conds = append(conds, fmt.Sprintf("p.category = $%d", i))
		args = append(args, f.Category)
		i++
	}
	if f.Supplier != "" {
		conds = append(conds, fmt.Sprintf("ib.supplier = $%d", i))
		args = append(args, f.Supplier)
		i++
	}
	if f.Search != "" {
		conds = append(conds, fmt.Sprintf("(p.name ILIKE $%d OR p.barcode = $%d OR ib.batch_number ILIKE $%d)", i, i+1, i+2))
		args = append(args, "%"+f.Search+"%", f.Search, "%"+f.Search+"%")
		i += 3
	}
	if f.ExpiryFrom != nil {
		conds = append(conds, fmt.Sprintf("ib.expiry_date >= $%d", i))
		args = append(args, *f.ExpiryFrom)
		i++
	}
	if f.ExpiryTo != nil {
		conds = append(conds, fmt.Sprintf("ib.expiry_date <= $%d", i))
		args = append(args, *f.ExpiryTo)
		i++
	}

	where := strings.Join(conds, " AND ")

	allowedSort := map[string]bool{
		"p.name": true, "ib.expiry_date": true, "ra.risk_level": true,
		"ra.estimated_loss": true, "ib.current_quantity": true,
	}
	sortCol := "ra.risk_level"
	if allowedSort[f.SortBy] {
		sortCol = f.SortBy
	}
	sortDir := "ASC"
	if strings.ToUpper(f.SortDir) == "DESC" {
		sortDir = "DESC"
	}

	countSQL := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM inventory_batches ib
		JOIN products p ON p.id = ib.product_id
		LEFT JOIN risk_assessments ra ON ra.batch_id = ib.id
		WHERE %s`, where)

	var total int
	if err := r.pool.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count inventory: %w", err)
	}

	offset := (f.Page - 1) * f.Limit
	args = append(args, f.Limit, offset)

	riskOrder := `CASE ra.risk_level WHEN 'CRITICAL' THEN 1 WHEN 'HIGH' THEN 2 WHEN 'MEDIUM' THEN 3 ELSE 4 END`
	orderExpr := fmt.Sprintf("%s %s", sortCol, sortDir)
	if sortCol == "ra.risk_level" {
		orderExpr = fmt.Sprintf("%s %s", riskOrder, sortDir)
	}

	dataSQL := fmt.Sprintf(`
		SELECT ib.id, ib.pharmacy_id, ib.product_id, p.name, COALESCE(p.name_el,''), COALESCE(p.category,''),
		       COALESCE(p.manufacturer,''), COALESCE(p.barcode,''), COALESCE(ib.batch_number,''),
		       ib.expiry_date, ib.initial_quantity, ib.current_quantity,
		       ib.purchase_price, ib.selling_price, COALESCE(ib.supplier,''), ib.received_date,
		       COALESCE(ra.risk_level,'LOW'), COALESCE(ra.days_until_expiry, (ib.expiry_date - CURRENT_DATE)::int),
		       COALESCE(ra.avg_daily_sales,0), COALESCE(ra.expected_sales,0),
		       COALESCE(ra.estimated_surplus,0), COALESCE(ra.estimated_loss,0),
		       COALESCE(ra.suggested_discount_percent,0)
		FROM inventory_batches ib
		JOIN products p ON p.id = ib.product_id
		LEFT JOIN risk_assessments ra ON ra.batch_id = ib.id
		WHERE %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d`, where, orderExpr, i, i+1)

	rows, err := r.pool.Query(ctx, dataSQL, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list inventory: %w", err)
	}
	defer rows.Close()

	var batches []models.InventoryBatchWithRisk
	for rows.Next() {
		var b models.InventoryBatchWithRisk
		if err := rows.Scan(
			&b.ID, &b.PharmacyID, &b.ProductID, &b.ProductName, &b.ProductNameEl, &b.Category,
			&b.Manufacturer, &b.Barcode, &b.BatchNumber,
			&b.ExpiryDate, &b.InitialQuantity, &b.CurrentQuantity,
			&b.PurchasePrice, &b.SellingPrice, &b.Supplier, &b.ReceivedDate,
			&b.RiskLevel, &b.DaysUntilExpiry, &b.AvgDailySales, &b.ExpectedSales,
			&b.EstimatedSurplus, &b.EstimatedLoss, &b.SuggestedDiscountPct,
		); err != nil {
			return nil, 0, err
		}
		batches = append(batches, b)
	}
	return batches, total, nil
}

func (r *Repository) GetBatch(ctx context.Context, id uuid.UUID) (*models.InventoryBatch, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, pharmacy_id, product_id, COALESCE(batch_number,''), expiry_date,
		       initial_quantity, current_quantity, purchase_price, selling_price,
		       COALESCE(supplier,''), received_date, created_at, updated_at
		FROM inventory_batches WHERE id = $1`, id)

	var b models.InventoryBatch
	err := row.Scan(&b.ID, &b.PharmacyID, &b.ProductID, &b.BatchNumber, &b.ExpiryDate,
		&b.InitialQuantity, &b.CurrentQuantity, &b.PurchasePrice, &b.SellingPrice,
		&b.Supplier, &b.ReceivedDate, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get batch: %w", err)
	}
	return &b, nil
}

func (r *Repository) GetBatchWithRisk(ctx context.Context, id uuid.UUID) (*models.InventoryBatchWithRisk, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT ib.id, ib.pharmacy_id, ib.product_id, p.name, COALESCE(p.name_el,''), COALESCE(p.category,''),
		       COALESCE(p.manufacturer,''), COALESCE(p.barcode,''), COALESCE(ib.batch_number,''),
		       ib.expiry_date, ib.initial_quantity, ib.current_quantity,
		       ib.purchase_price, ib.selling_price, COALESCE(ib.supplier,''), ib.received_date,
		       COALESCE(ra.risk_level,'LOW'), COALESCE(ra.days_until_expiry, (ib.expiry_date - CURRENT_DATE)::int),
		       COALESCE(ra.avg_daily_sales,0), COALESCE(ra.expected_sales,0),
		       COALESCE(ra.estimated_surplus,0), COALESCE(ra.estimated_loss,0),
		       COALESCE(ra.suggested_discount_percent,0)
		FROM inventory_batches ib
		JOIN products p ON p.id = ib.product_id
		LEFT JOIN risk_assessments ra ON ra.batch_id = ib.id
		WHERE ib.id = $1`, id)

	var b models.InventoryBatchWithRisk
	err := row.Scan(
		&b.ID, &b.PharmacyID, &b.ProductID, &b.ProductName, &b.ProductNameEl, &b.Category,
		&b.Manufacturer, &b.Barcode, &b.BatchNumber,
		&b.ExpiryDate, &b.InitialQuantity, &b.CurrentQuantity,
		&b.PurchasePrice, &b.SellingPrice, &b.Supplier, &b.ReceivedDate,
		&b.RiskLevel, &b.DaysUntilExpiry, &b.AvgDailySales, &b.ExpectedSales,
		&b.EstimatedSurplus, &b.EstimatedLoss, &b.SuggestedDiscountPct,
	)
	if err != nil {
		return nil, fmt.Errorf("get batch with risk: %w", err)
	}
	return &b, nil
}

type BatchImportRow struct {
	ProductID     uuid.UUID
	PharmacyID    uuid.UUID
	BatchNumber   string
	ExpiryDate    time.Time
	Quantity      int
	PurchasePrice float64
	SellingPrice  float64
	Supplier      string
}

func (r *Repository) InsertBatch(ctx context.Context, row BatchImportRow) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.pool.QueryRow(ctx, `
		INSERT INTO inventory_batches (
			id, pharmacy_id, product_id, batch_number, expiry_date,
			initial_quantity, current_quantity, purchase_price, selling_price,
			supplier, received_date
		) VALUES (
			gen_random_uuid(), $1, $2, NULLIF($3,''), $4, $5, $5, $6, $7, NULLIF($8,''), CURRENT_DATE
		) RETURNING id`,
		row.PharmacyID, row.ProductID, row.BatchNumber, row.ExpiryDate,
		row.Quantity, row.PurchasePrice, row.SellingPrice, row.Supplier,
	).Scan(&id)
	return id, err
}

func (r *Repository) ListSuppliers(ctx context.Context, pharmacyID uuid.UUID) ([]string, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT DISTINCT supplier FROM inventory_batches
		WHERE pharmacy_id = $1 AND supplier IS NOT NULL AND supplier != ''
		ORDER BY supplier`, pharmacyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suppliers []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		suppliers = append(suppliers, s)
	}
	return suppliers, nil
}
