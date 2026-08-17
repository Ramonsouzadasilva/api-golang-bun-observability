package logger

import (
	"log/slog"
	"os"
)

// New cria um logger estruturado. Em produção usa JSON (fácil de indexar
// em ferramentas de observabilidade); em desenvolvimento usa texto legível.
func New(env string) *slog.Logger {
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}

	var handler slog.Handler
	if env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		opts.Level = slog.LevelDebug
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	log := slog.New(handler)
	slog.SetDefault(log)

	return log
}
