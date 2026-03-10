package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"finapp/internal/config"
	"finapp/internal/db"
	"finapp/internal/handler"
	"finapp/internal/model"
	"finapp/internal/pluggy"
	"finapp/internal/repository"
	"finapp/internal/service"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	cfg, err := config.Load()
	if err != nil {
		slog.Error("load config", "err", err)
		os.Exit(1)
	}

	ctx := context.Background()

	// Database pool
	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("connect to database", "err", err)
		os.Exit(1)
	}
	defer pool.Close()
	slog.Info("database connected")

	// Run migrations
	if err := runMigrations(cfg.DatabaseURL); err != nil {
		slog.Error("run migrations", "err", err)
		os.Exit(1)
	}
	slog.Info("migrations applied")

	// Repositories
	userRepo := repository.NewUserRepository(pool)
	pluggyItemRepo := repository.NewPluggyItemRepository(pool)
	accountRepo := repository.NewAccountRepository(pool)
	txRepo := repository.NewTransactionRepository(pool)
	categoryRepo := repository.NewCategoryRepository(pool)
	budgetRepo := repository.NewBudgetRepository(pool)
	goalRepo := repository.NewGoalRepository(pool)
	webhookLogRepo := repository.NewWebhookLogRepository(pool)

	// Seed system categories
	if err := seedCategories(ctx, categoryRepo); err != nil {
		slog.Error("seed categories", "err", err)
		os.Exit(1)
	}

	// Pluggy client
	pluggyClient := pluggy.NewClient(cfg.PluggyBaseURL, cfg.PluggyClientID, cfg.PluggyClientSecret)

	// Services
	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTAccessTTLMinutes, cfg.JWTRefreshTTLDays)
	pluggySyncSvc := service.NewPluggySyncService(pluggyClient, pluggyItemRepo, accountRepo, txRepo, webhookLogRepo)
	accountSvc := service.NewAccountService(accountRepo)
	txSvc := service.NewTransactionService(txRepo)
	categorySvc := service.NewCategoryService(categoryRepo)
	budgetSvc := service.NewBudgetService(budgetRepo, txRepo)
	goalSvc := service.NewGoalService(goalRepo)
	reportSvc := service.NewReportService(txRepo)
	simulationSvc := service.NewSimulationService()
	projectionSvc := service.NewProjectionService(txRepo, accountRepo)

	// Handlers
	handlers := handler.Handlers{
		Auth:        handler.NewAuthHandler(authSvc),
		Pluggy:      handler.NewPluggyHandler(pluggySyncSvc, cfg.PluggyWebhookSecret),
		Account:     handler.NewAccountHandler(accountSvc),
		Transaction: handler.NewTransactionHandler(txSvc),
		Category:    handler.NewCategoryHandler(categorySvc),
		Budget:      handler.NewBudgetHandler(budgetSvc),
		Goal:        handler.NewGoalHandler(goalSvc),
		Report:      handler.NewReportHandler(reportSvc),
		Simulation:  handler.NewSimulationHandler(simulationSvc),
		Projection:  handler.NewProjectionHandler(projectionSvc),
	}

	router := handler.NewRouter(handlers, authSvc)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("server starting", "port", cfg.Port, "env", cfg.Env)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	<-quit
	slog.Info("shutting down gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("forced shutdown", "err", err)
	}
	slog.Info("server stopped")
}

func runMigrations(databaseURL string) error {
	m, err := migrate.New("file://migrations", databaseURL)
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}

func seedCategories(ctx context.Context, repo repository.CategoryRepository) error {
	has, err := repo.HasSystemCategories(ctx)
	if err != nil {
		return err
	}
	if has {
		return nil
	}
	return repo.SeedSystemCategories(ctx, model.DefaultSystemCategories)
}
