package handlers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"pharmasense/backend/internal/middleware"
	"pharmasense/backend/internal/repository"
	"pharmasense/backend/internal/services"
)

type InventoryHandler struct {
	repo   *repository.Repository
	engine *services.RiskEngine
}

func NewInventoryHandler(repo *repository.Repository, engine *services.RiskEngine) *InventoryHandler {
	return &InventoryHandler{repo: repo, engine: engine}
}

func (h *InventoryHandler) List(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)

	f := repository.InventoryFilter{
		PharmacyID: claims.PharmacyID,
		RiskLevel:  queryStr(r, "risk_level"),
		Category:   queryStr(r, "category"),
		Supplier:   queryStr(r, "supplier"),
		Search:     queryStr(r, "search"),
		Page:       queryInt(r, "page", 1),
		Limit:      queryInt(r, "limit", 50),
		SortBy:     queryStr(r, "sort_by"),
		SortDir:    queryStr(r, "sort_dir"),
	}

	if from := queryStr(r, "expiry_from"); from != "" {
		if t, err := time.Parse("2006-01-02", from); err == nil {
			f.ExpiryFrom = &t
		}
	}
	if to := queryStr(r, "expiry_to"); to != "" {
		if t, err := time.Parse("2006-01-02", to); err == nil {
			f.ExpiryTo = &t
		}
	}

	batches, total, err := h.repo.ListInventoryWithRisk(r.Context(), f)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "list inventory failed")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"data":  batches,
		"total": total,
		"page":  f.Page,
		"limit": f.Limit,
	})
}

func (h *InventoryHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	b, err := h.repo.GetBatchWithRisk(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "batch not found")
		return
	}
	writeJSON(w, http.StatusOK, b)
}

func (h *InventoryHandler) GetFilters(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)

	cats, _ := h.repo.ListCategories(r.Context())
	suppliers, _ := h.repo.ListSuppliers(r.Context(), claims.PharmacyID)

	writeJSON(w, http.StatusOK, map[string]any{
		"categories": cats,
		"suppliers":  suppliers,
		"risk_levels": []string{"CRITICAL", "HIGH", "MEDIUM", "LOW"},
	})
}

func (h *InventoryHandler) ExportCSV(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)

	f := repository.InventoryFilter{
		PharmacyID: claims.PharmacyID,
		Page:       1,
		Limit:      10000,
	}

	batches, _, err := h.repo.ListInventoryWithRisk(r.Context(), f)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "export failed")
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", `attachment; filename="inventory_export.csv"`)

	cw := csv.NewWriter(w)
	_ = cw.Write([]string{
		"Product", "Category", "Batch#", "Expiry Date", "Days Until Expiry",
		"Quantity", "Risk Level", "Estimated Loss (€)", "Supplier",
	})

	for _, b := range batches {
		_ = cw.Write([]string{
			b.ProductName,
			b.Category,
			b.BatchNumber,
			b.ExpiryDate.Format("2006-01-02"),
			strconv.Itoa(b.DaysUntilExpiry),
			strconv.Itoa(b.CurrentQuantity),
			string(b.RiskLevel),
			fmt.Sprintf("%.2f", b.EstimatedLoss),
			b.Supplier,
		})
	}
	cw.Flush()
}

func (h *InventoryHandler) ImportCSV(w http.ResponseWriter, r *http.Request) {
	// Parse CSV import — validate and preview
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "file too large")
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "no file uploaded")
		return
	}
	defer file.Close()

	cr := csv.NewReader(file)
	records, err := cr.ReadAll()
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid CSV format")
		return
	}

	if len(records) < 2 {
		writeError(w, http.StatusBadRequest, "CSV must have header + at least one row")
		return
	}

	header := records[0]
	requiredCols := []string{"product_name", "barcode", "batch_number", "expiry_date", "quantity", "purchase_price"}
	colIndex := make(map[string]int)
	for i, h := range header {
		colIndex[strings.ToLower(strings.TrimSpace(h))] = i
	}
	var missingCols []string
	for _, col := range requiredCols {
		if _, ok := colIndex[col]; !ok {
			missingCols = append(missingCols, col)
		}
	}
	if len(missingCols) > 0 {
		writeError(w, http.StatusBadRequest, "missing columns: "+strings.Join(missingCols, ", "))
		return
	}

	type PreviewRow struct {
		Row          int    `json:"row"`
		ProductName  string `json:"product_name"`
		Barcode      string `json:"barcode"`
		BatchNumber  string `json:"batch_number"`
		ExpiryDate   string `json:"expiry_date"`
		Quantity     string `json:"quantity"`
		PurchasePrice string `json:"purchase_price"`
		Error        string `json:"error,omitempty"`
	}

	var preview []PreviewRow
	var errCount int
	for i, rec := range records[1:] {
		if i >= 20 { // preview first 20 rows
			break
		}
		pr := PreviewRow{
			Row:          i + 2,
			ProductName:  rec[colIndex["product_name"]],
			Barcode:      rec[colIndex["barcode"]],
			BatchNumber:  rec[colIndex["batch_number"]],
			ExpiryDate:   rec[colIndex["expiry_date"]],
			Quantity:     rec[colIndex["quantity"]],
			PurchasePrice: rec[colIndex["purchase_price"]],
		}
		// Validate
		if _, err := time.Parse("2006-01-02", pr.ExpiryDate); err != nil {
			pr.Error = "invalid expiry_date format (YYYY-MM-DD expected)"
			errCount++
		}
		preview = append(preview, pr)
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"total_rows":   len(records) - 1,
		"preview":      preview,
		"error_count":  errCount,
		"template_url": "/api/v1/inventory/import/template",
	})
}

func (h *InventoryHandler) ImportTemplate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", `attachment; filename="import_template.csv"`)

	cw := csv.NewWriter(w)
	_ = cw.Write([]string{"product_name", "barcode", "batch_number", "expiry_date", "quantity", "purchase_price", "selling_price", "supplier"})
	_ = cw.Write([]string{"Paracetamol 500mg", "5201234567890", "BTH-2024-00001", "2026-06-01", "100", "1.50", "2.20", "MedSupply Cyprus"})
	_ = cw.Write([]string{"Vitamin C 1000mg", "5201234567891", "BTH-2024-00002", "2026-12-01", "200", "3.00", "4.50", "EuroMeds"})
	cw.Flush()
}

func (h *InventoryHandler) ConfirmImport(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "file too large")
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "no file uploaded")
		return
	}
	defer file.Close()

	cr := csv.NewReader(file)
	records, err := cr.ReadAll()
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid CSV format")
		return
	}
	if len(records) < 2 {
		writeError(w, http.StatusBadRequest, "CSV must have header + at least one row")
		return
	}

	header := records[0]
	colIndex := make(map[string]int)
	for i, h := range header {
		colIndex[strings.ToLower(strings.TrimSpace(h))] = i
	}
	required := []string{"product_name", "barcode", "batch_number", "expiry_date", "quantity", "purchase_price"}
	var missing []string
	for _, col := range required {
		if _, ok := colIndex[col]; !ok {
			missing = append(missing, col)
		}
	}
	if len(missing) > 0 {
		writeError(w, http.StatusBadRequest, "missing columns: "+strings.Join(missing, ", "))
		return
	}

	type rowErr struct {
		Row   int    `json:"row"`
		Error string `json:"error"`
	}
	var errs []rowErr
	imported := 0

	today := time.Now().UTC().Truncate(24 * time.Hour)

	for i, rec := range records[1:] {
		rowNum := i + 2
		get := func(col string) string {
			idx, ok := colIndex[col]
			if !ok || idx >= len(rec) {
				return ""
			}
			return strings.TrimSpace(rec[idx])
		}

		productName := get("product_name")
		barcode := get("barcode")
		batchNumber := get("batch_number")
		expiryStr := get("expiry_date")
		qtyStr := get("quantity")
		priceStr := get("purchase_price")

		expiryDate, err := time.Parse("2006-01-02", expiryStr)
		if err != nil {
			errs = append(errs, rowErr{Row: rowNum, Error: "invalid expiry_date (YYYY-MM-DD expected)"})
			continue
		}
		qty, err := strconv.Atoi(qtyStr)
		if err != nil || qty <= 0 {
			errs = append(errs, rowErr{Row: rowNum, Error: "invalid quantity"})
			continue
		}
		purchasePrice, err := strconv.ParseFloat(priceStr, 64)
		if err != nil || purchasePrice <= 0 {
			errs = append(errs, rowErr{Row: rowNum, Error: "invalid purchase_price"})
			continue
		}

		sellingPrice := purchasePrice * 1.5
		if sp, err := strconv.ParseFloat(get("selling_price"), 64); err == nil && sp > 0 {
			sellingPrice = sp
		}

		productID, err := h.repo.FindOrCreateProduct(r.Context(), productName, barcode)
		if err != nil {
			errs = append(errs, rowErr{Row: rowNum, Error: "product lookup failed: " + err.Error()})
			continue
		}

		batchID, err := h.repo.InsertBatch(r.Context(), repository.BatchImportRow{
			ProductID:     productID,
			PharmacyID:    claims.PharmacyID,
			BatchNumber:   batchNumber,
			ExpiryDate:    expiryDate,
			Quantity:      qty,
			PurchasePrice: purchasePrice,
			SellingPrice:  sellingPrice,
			Supplier:      get("supplier"),
		})
		if err != nil {
			errs = append(errs, rowErr{Row: rowNum, Error: "insert failed: " + err.Error()})
			continue
		}

		avg, _ := h.repo.AvgDailySales(r.Context(), claims.PharmacyID, productID)
		days := int(expiryDate.Sub(today).Hours() / 24)
		ra := services.Calculate(days, qty, avg, purchasePrice)
		ra.BatchID = batchID
		ra.PharmacyID = claims.PharmacyID
		_ = h.repo.UpsertRiskAssessment(r.Context(), &ra)

		imported++
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"imported": imported,
		"errors":   errs,
	})
}
