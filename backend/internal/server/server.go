package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"pharmasense/backend/internal/handlers"
	"pharmasense/backend/internal/middleware"
	"pharmasense/backend/internal/repository"
	"pharmasense/backend/internal/services"
)

func New(repo *repository.Repository, authSvc *services.AuthService, engine *services.RiskEngine) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.Recoverer)
	r.Use(chimw.RealIP)
	r.Use(middleware.CORS)

	authH := handlers.NewAuthHandler(authSvc, repo)
	pharmH := handlers.NewPharmacyHandler(repo)
	invH := handlers.NewInventoryHandler(repo, engine)
	riskH := handlers.NewRiskHandler(repo, engine)
	alertsH := handlers.NewAlertsHandler(repo)
	reportsH := handlers.NewReportsHandler(repo)
	productsH := handlers.NewProductsHandler(repo)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":"ok"}`))
		})

		// Public auth routes
		r.Post("/auth/login", authH.Login)
		r.Post("/auth/logout", authH.Logout)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authSvc))

			r.Get("/auth/me", authH.Me)

			r.Get("/pharmacy", pharmH.Get)
			r.Patch("/pharmacy", pharmH.Update)
			r.Get("/pharmacy/users", pharmH.GetUsers)

			r.Get("/products", productsH.List)
			r.Get("/products/{id}", productsH.Get)

			r.Get("/inventory", invH.List)
			r.Get("/inventory/filters", invH.GetFilters)
			r.Get("/inventory/export", invH.ExportCSV)
			r.Post("/inventory/import", invH.ImportCSV)
			r.Post("/inventory/import/confirm", invH.ConfirmImport)
			r.Get("/inventory/import/template", invH.ImportTemplate)
			r.Get("/inventory/{id}", invH.Get)

			r.Get("/risk/dashboard", riskH.Dashboard)
			r.Get("/risk/timeline", riskH.Timeline)
			r.Get("/risk/top", riskH.TopRisk)
			r.Get("/risk/assessments", riskH.Assessments)
			r.Post("/risk/recalculate", riskH.Recalculate)

			r.Get("/alerts", alertsH.List)
			r.Post("/alerts/{batch_id}/action", alertsH.TakeAction)
			r.Get("/alerts/history", alertsH.History)

			r.Get("/reports/savings", reportsH.Savings)
			r.Get("/reports/categories", reportsH.Categories)
			r.Get("/reports/waste", reportsH.Waste)
		})
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	return r
}
