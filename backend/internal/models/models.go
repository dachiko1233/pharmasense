package models

import (
	"time"

	"github.com/google/uuid"
)

type Pharmacy struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	LicenseNumber string    `json:"license_number"`
	Address       string    `json:"address"`
	City          string    `json:"city"`
	Phone         string    `json:"phone"`
	Email         string    `json:"email"`
	Language      string    `json:"language"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type User struct {
	ID          uuid.UUID  `json:"id"`
	PharmacyID  uuid.UUID  `json:"pharmacy_id"`
	Email       string     `json:"email"`
	PasswordHash string    `json:"-"`
	FullName    string     `json:"full_name"`
	Role        string     `json:"role"`
	IsActive    bool       `json:"is_active"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

type Product struct {
	ID                   uuid.UUID `json:"id"`
	Barcode              string    `json:"barcode"`
	Name                 string    `json:"name"`
	NameEl               string    `json:"name_el"`
	Category             string    `json:"category"`
	Manufacturer         string    `json:"manufacturer"`
	RequiresPrescription bool      `json:"requires_prescription"`
	CreatedAt            time.Time `json:"created_at"`
}

type InventoryBatch struct {
	ID              uuid.UUID `json:"id"`
	PharmacyID      uuid.UUID `json:"pharmacy_id"`
	ProductID       uuid.UUID `json:"product_id"`
	BatchNumber     string    `json:"batch_number"`
	ExpiryDate      time.Time `json:"expiry_date"`
	InitialQuantity int       `json:"initial_quantity"`
	CurrentQuantity int       `json:"current_quantity"`
	PurchasePrice   float64   `json:"purchase_price"`
	SellingPrice    float64   `json:"selling_price"`
	Supplier        string    `json:"supplier"`
	ReceivedDate    time.Time `json:"received_date"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Sale struct {
	ID          uuid.UUID `json:"id"`
	PharmacyID  uuid.UUID `json:"pharmacy_id"`
	BatchID     uuid.UUID `json:"batch_id"`
	ProductID   uuid.UUID `json:"product_id"`
	Quantity    int       `json:"quantity"`
	UnitPrice   float64   `json:"unit_price"`
	TotalAmount float64   `json:"total_amount"`
	SaleDate    time.Time `json:"sale_date"`
	CreatedAt   time.Time `json:"created_at"`
}

type RiskLevel string

const (
	RiskCritical RiskLevel = "CRITICAL"
	RiskHigh     RiskLevel = "HIGH"
	RiskMedium   RiskLevel = "MEDIUM"
	RiskLow      RiskLevel = "LOW"
)

type RiskAssessment struct {
	ID                    uuid.UUID `json:"id"`
	BatchID               uuid.UUID `json:"batch_id"`
	PharmacyID            uuid.UUID `json:"pharmacy_id"`
	RiskLevel             RiskLevel `json:"risk_level"`
	DaysUntilExpiry       int       `json:"days_until_expiry"`
	AvgDailySales         float64   `json:"avg_daily_sales"`
	ExpectedSales         int       `json:"expected_sales"`
	EstimatedSurplus      int       `json:"estimated_surplus"`
	EstimatedLoss         float64   `json:"estimated_loss"`
	SuggestedDiscountPct  int       `json:"suggested_discount_percent"`
	CalculatedAt          time.Time `json:"calculated_at"`
}

type AlertAction struct {
	ID              uuid.UUID `json:"id"`
	BatchID         uuid.UUID `json:"batch_id"`
	PharmacyID      uuid.UUID `json:"pharmacy_id"`
	UserID          uuid.UUID `json:"user_id"`
	ActionType      string    `json:"action_type"`
	DiscountPercent *int      `json:"discount_percent"`
	Notes           string    `json:"notes"`
	CreatedAt       time.Time `json:"created_at"`
}

// InventoryBatchWithRisk combines batch + product + risk for API responses.
type InventoryBatchWithRisk struct {
	ID                   uuid.UUID `json:"id"`
	PharmacyID           uuid.UUID `json:"pharmacy_id"`
	ProductID            uuid.UUID `json:"product_id"`
	ProductName          string    `json:"product_name"`
	ProductNameEl        string    `json:"product_name_el"`
	Category             string    `json:"category"`
	Manufacturer         string    `json:"manufacturer"`
	Barcode              string    `json:"barcode"`
	BatchNumber          string    `json:"batch_number"`
	ExpiryDate           time.Time `json:"expiry_date"`
	InitialQuantity      int       `json:"initial_quantity"`
	CurrentQuantity      int       `json:"current_quantity"`
	PurchasePrice        float64   `json:"purchase_price"`
	SellingPrice         float64   `json:"selling_price"`
	Supplier             string    `json:"supplier"`
	ReceivedDate         time.Time `json:"received_date"`
	RiskLevel            RiskLevel `json:"risk_level"`
	DaysUntilExpiry      int       `json:"days_until_expiry"`
	AvgDailySales        float64   `json:"avg_daily_sales"`
	ExpectedSales        int       `json:"expected_sales"`
	EstimatedSurplus     int       `json:"estimated_surplus"`
	EstimatedLoss        float64   `json:"estimated_loss"`
	SuggestedDiscountPct int       `json:"suggested_discount_percent"`
}

type DashboardData struct {
	CriticalCount        int                    `json:"critical_count"`
	HighCount            int                    `json:"high_count"`
	MediumCount          int                    `json:"medium_count"`
	LowCount             int                    `json:"low_count"`
	TotalEstimatedLoss   float64                `json:"total_estimated_loss"`
	TotalPotentialSavings float64               `json:"total_potential_savings"`
	TotalInventoryValue  float64                `json:"total_inventory_value"`
	CriticalAlerts       []InventoryBatchWithRisk `json:"critical_alerts"`
}

type TimelineMonth struct {
	Month    string `json:"month"`
	Critical int    `json:"critical"`
	High     int    `json:"high"`
	Medium   int    `json:"medium"`
	Low      int    `json:"low"`
}

type TopRiskItem struct {
	ProductName   string    `json:"product_name"`
	BatchNumber   string    `json:"batch_number"`
	EstimatedLoss float64   `json:"estimated_loss"`
	RiskLevel     RiskLevel `json:"risk_level"`
	Category      string    `json:"category"`
}

type AlertActionWithDetails struct {
	AlertAction
	UserName    string `json:"user_name"`
	ProductName string `json:"product_name"`
	BatchNumber string `json:"batch_number"`
}

type SavingsReport struct {
	Month        string  `json:"month"`
	ActionsCount int     `json:"actions_count"`
	EstimatedSaved float64 `json:"estimated_saved"`
}

type CategoryReport struct {
	Category      string  `json:"category"`
	TotalBatches  int     `json:"total_batches"`
	AtRisk        int     `json:"at_risk"`
	EstimatedLoss float64 `json:"estimated_loss"`
}

type ListResponse[T any] struct {
	Data  []T `json:"data"`
	Total int `json:"total"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
}
