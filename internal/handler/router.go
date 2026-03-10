package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	appMiddleware "finapp/internal/middleware"
	"finapp/internal/service"
)

type Handlers struct {
	Auth        *AuthHandler
	Pluggy      *PluggyHandler
	Account     *AccountHandler
	Transaction *TransactionHandler
	Category    *CategoryHandler
	Budget      *BudgetHandler
	Goal        *GoalHandler
	Report      *ReportHandler
	Simulation  *SimulationHandler
	Projection  *ProjectionHandler
}

func NewRouter(h Handlers, authSvc service.AuthService) http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chiMiddleware.RequestID)
	r.Use(appMiddleware.Recovery)
	r.Use(appMiddleware.Logging)
	r.Use(chiMiddleware.CleanPath)

	r.Route("/api/v1", func(r chi.Router) {
		// Public: auth
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", h.Auth.Register)
			r.Post("/login", h.Auth.Login)
			r.Post("/refresh", h.Auth.Refresh)
		})

		// Public: Pluggy webhook (no JWT — HMAC verified inside handler)
		r.Post("/pluggy/webhook", h.Pluggy.Webhook)

		// Protected: all other routes require JWT
		r.Group(func(r chi.Router) {
			r.Use(appMiddleware.Auth(authSvc))

			// Pluggy
			r.Route("/pluggy", func(r chi.Router) {
				r.Post("/connect-token", h.Pluggy.GenerateConnectToken)
				r.Get("/items", h.Pluggy.ListItems)
				r.Delete("/items/{id}", h.Pluggy.DisconnectItem)
				r.Post("/sync", h.Pluggy.Sync)
			})

			// Accounts
			r.Route("/accounts", func(r chi.Router) {
				r.Get("/", h.Account.List)
				r.Get("/{id}", h.Account.Get)
			})

			// Transactions
			r.Route("/transactions", func(r chi.Router) {
				r.Get("/", h.Transaction.List)
				r.Get("/{id}", h.Transaction.Get)
				r.Patch("/{id}", h.Transaction.Update)
			})

			// Categories
			r.Route("/categories", func(r chi.Router) {
				r.Get("/", h.Category.List)
				r.Post("/", h.Category.Create)
				r.Get("/{id}", h.Category.Get)
				r.Put("/{id}", h.Category.Update)
				r.Delete("/{id}", h.Category.Delete)
			})

			// Budgets
			r.Route("/budgets", func(r chi.Router) {
				r.Get("/", h.Budget.List)
				r.Post("/", h.Budget.Create)
				r.Get("/{id}", h.Budget.Get)
				r.Put("/{id}", h.Budget.Update)
				r.Delete("/{id}", h.Budget.Delete)
				r.Get("/{id}/progress", h.Budget.Progress)
			})

			// Goals
			r.Route("/goals", func(r chi.Router) {
				r.Get("/", h.Goal.List)
				r.Post("/", h.Goal.Create)
				r.Get("/{id}", h.Goal.Get)
				r.Put("/{id}", h.Goal.Update)
				r.Delete("/{id}", h.Goal.Delete)
				r.Get("/{id}/progress", h.Goal.Progress)
			})

			// Reports
			r.Route("/reports", func(r chi.Router) {
				r.Get("/summary", h.Report.Summary)
				r.Get("/by-category", h.Report.ByCategory)
				r.Get("/cashflow", h.Report.CashFlow)
			})

			// Simulations (no DB needed)
			r.Route("/simulations", func(r chi.Router) {
				r.Post("/compound-interest", h.Simulation.CompoundInterest)
				r.Post("/loan", h.Simulation.Loan)
			})

			// Projections
			r.Route("/projections", func(r chi.Router) {
				r.Get("/balance", h.Projection.Balance)
			})
		})
	})

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	return r
}
