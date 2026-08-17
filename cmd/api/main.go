package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"

	"go-api-base/internal/config"
	"go-api-base/internal/database"
	"go-api-base/internal/database/migrations"
	"go-api-base/internal/database/seeds"
	"go-api-base/internal/logger"
	authmodule "go-api-base/internal/modules/auth"
	usermodule "go-api-base/internal/modules/user"
	"go-api-base/internal/observability"
	"go-api-base/internal/router"
	"go-api-base/pkg/hash"
	"go-api-base/pkg/jwt"
)

func main() {
	observability.Register()
	migrateFlag := flag.String("migrate", "", "run migrations: up | down | status")
	seedFlag := flag.Bool("seed", false, "seed the database with initial data")
	flag.Parse()

	cfg := config.Load()
	log := logger.New(cfg.App.Env)

	hash.Cost = cfg.App.BcryptCost

	db, err := database.Connect(cfg.DB, cfg.App.Env == "development")
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if *migrateFlag != "" {
		runMigrations(db, *migrateFlag, log)
		return
	}

	if *seedFlag {
		if err := seeds.Run(context.Background(), db); err != nil {
			log.Error("failed to seed database", "error", err)
			os.Exit(1)
		}
		return
	}

	startServer(cfg, db, log)
}

func startServer(cfg *config.Config, db *bun.DB, log *slog.Logger) {
	jwtManager := jwt.NewManager(cfg.JWT.Secret)

	userMod := usermodule.NewModule(db)
	authMod := authmodule.NewModule(db, userMod.Repository, jwtManager, cfg.JWT)

	handler := router.New(router.Dependencies{
		Logger:     log,
		JWTManager: jwtManager,
		AuthModule: authMod,
		UserModule: userMod,
	})

	srv := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info("server started", "port", cfg.App.Port, "env", cfg.App.Env)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	go func() {
		log.Info("pprof server started", "port", 6060)

		// Habilita profiling de mutex e block (útil para achar gargalos em concorrência)
		runtime.SetMutexProfileFraction(1)
		runtime.SetBlockProfileRate(1)

		mux := http.NewServeMux()
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

		if err := http.ListenAndServe(":6060", mux); err != nil {
			log.Error("pprof server error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("forced shutdown", "error", err)
	}

	log.Info("server stopped gracefully")
}

func runMigrations(db *bun.DB, action string, log *slog.Logger) {
	migrator := migrate.NewMigrator(db, migrations.Migrations)
	ctx := context.Background()

	if err := migrator.Init(ctx); err != nil {
		log.Error("failed to init migrator", "error", err)
		os.Exit(1)
	}

	switch action {
	case "up":
		group, err := migrator.Migrate(ctx)
		if err != nil {
			log.Error("migration failed", "error", err)
			os.Exit(1)
		}
		if group.IsZero() {
			log.Info("no new migrations to run")
			return
		}
		log.Info("migrated successfully", "group", group.String())

	case "down":
		group, err := migrator.Rollback(ctx)
		if err != nil {
			log.Error("rollback failed", "error", err)
			os.Exit(1)
		}
		log.Info("rolled back successfully", "group", group.String())

	case "status":
		ms, err := migrator.MigrationsWithStatus(ctx)
		if err != nil {
			log.Error("failed to get migration status", "error", err)
			os.Exit(1)
		}
		fmt.Println(ms)

	default:
		log.Error("unknown migrate action, use up | down | status", "action", action)
		os.Exit(1)
	}
}
