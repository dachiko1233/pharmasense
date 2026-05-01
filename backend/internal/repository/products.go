package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"pharmasense/backend/internal/models"
)

type ProductFilter struct {
	Search   string
	Category string
	Page     int
	Limit    int
}

func (r *Repository) ListProducts(ctx context.Context, f ProductFilter) ([]models.Product, int, error) {
	if f.Limit == 0 {
		f.Limit = 50
	}
	if f.Page == 0 {
		f.Page = 1
	}

	var conds []string
	var args []any
	i := 1

	if f.Search != "" {
		conds = append(conds, fmt.Sprintf("(name ILIKE $%d OR barcode = $%d)", i, i+1))
		args = append(args, "%"+f.Search+"%", f.Search)
		i += 2
	}
	if f.Category != "" {
		conds = append(conds, fmt.Sprintf("category = $%d", i))
		args = append(args, f.Category)
		i++
	}

	where := ""
	if len(conds) > 0 {
		where = "WHERE " + strings.Join(conds, " AND ")
	}

	var total int
	err := r.pool.QueryRow(ctx, fmt.Sprintf("SELECT COUNT(*) FROM products %s", where), args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (f.Page - 1) * f.Limit
	args = append(args, f.Limit, offset)
	rows, err := r.pool.Query(ctx, fmt.Sprintf(`
		SELECT id, COALESCE(barcode,''), name, COALESCE(name_el,''), COALESCE(category,''),
		       COALESCE(manufacturer,''), requires_prescription, created_at
		FROM products %s ORDER BY name LIMIT $%d OFFSET $%d`, where, i, i+1), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Barcode, &p.Name, &p.NameEl, &p.Category,
			&p.Manufacturer, &p.RequiresPrescription, &p.CreatedAt); err != nil {
			return nil, 0, err
		}
		products = append(products, p)
	}
	return products, total, nil
}

func (r *Repository) GetProduct(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, COALESCE(barcode,''), name, COALESCE(name_el,''), COALESCE(category,''),
		       COALESCE(manufacturer,''), requires_prescription, created_at
		FROM products WHERE id = $1`, id)

	var p models.Product
	err := row.Scan(&p.ID, &p.Barcode, &p.Name, &p.NameEl, &p.Category,
		&p.Manufacturer, &p.RequiresPrescription, &p.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get product: %w", err)
	}
	return &p, nil
}

func (r *Repository) FindOrCreateProduct(ctx context.Context, name, barcode string) (uuid.UUID, error) {
	if barcode != "" {
		var id uuid.UUID
		err := r.pool.QueryRow(ctx, `SELECT id FROM products WHERE barcode = $1`, barcode).Scan(&id)
		if err == nil {
			return id, nil
		}
	}
	var id uuid.UUID
	err := r.pool.QueryRow(ctx, `
		INSERT INTO products (id, barcode, name)
		VALUES (gen_random_uuid(), NULLIF($1,''), $2)
		ON CONFLICT (barcode) DO UPDATE SET name = EXCLUDED.name
		RETURNING id`, barcode, name).Scan(&id)
	return id, err
}

func (r *Repository) ListCategories(ctx context.Context) ([]string, error) {
	rows, err := r.pool.Query(ctx, `SELECT DISTINCT category FROM products WHERE category IS NOT NULL ORDER BY category`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []string
	for rows.Next() {
		var c string
		if err := rows.Scan(&c); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, nil
}
