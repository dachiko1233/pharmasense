package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"pharmasense/backend/internal/models"
)

func (r *Repository) GetPharmacy(ctx context.Context, id uuid.UUID) (*models.Pharmacy, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, name, license_number, COALESCE(address,''), COALESCE(city,''),
		       COALESCE(phone,''), COALESCE(email,''), language, created_at, updated_at
		FROM pharmacies WHERE id = $1`, id)

	var p models.Pharmacy
	err := row.Scan(&p.ID, &p.Name, &p.LicenseNumber, &p.Address, &p.City,
		&p.Phone, &p.Email, &p.Language, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get pharmacy: %w", err)
	}
	return &p, nil
}

func (r *Repository) UpdatePharmacy(ctx context.Context, p *models.Pharmacy) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE pharmacies SET name=$1, address=$2, city=$3, phone=$4, email=$5, language=$6, updated_at=NOW()
		WHERE id=$7`,
		p.Name, p.Address, p.City, p.Phone, p.Email, p.Language, p.ID)
	return err
}
