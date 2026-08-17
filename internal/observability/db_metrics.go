package observability

import (
	"context"
	"database/sql"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	DBPoolOpenConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_pool_open_connections",
			Help: "O número atual de conexões abertas com o banco de dados.",
		},
	)

	DBPoolIdleConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_pool_idle_connections",
			Help: "O número atual de conexões ociosas no pool.",
		},
	)

	DBPoolWaitCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_pool_wait_count_total",
			Help: "O número total de vezes que as conexões tiveram que esperar.",
		},
	)

	DBPoolWaitDuration = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_pool_wait_duration_seconds_total",
			Help: "O tempo total bloqueado esperando por uma nova conexão.",
		},
	)
)

func init() {
	prometheus.MustRegister(
		DBPoolOpenConnections,
		DBPoolIdleConnections,
		DBPoolWaitCount,
		DBPoolWaitDuration,
	)
}

// StartDBStatsCollector inicia uma goroutine que coleta periodicamente as estatísticas do banco de dados.
func StartDBStatsCollector(ctx context.Context, db *sql.DB, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				stats := db.Stats()
				DBPoolOpenConnections.Set(float64(stats.OpenConnections))
				DBPoolIdleConnections.Set(float64(stats.Idle))
				DBPoolWaitCount.Set(float64(stats.WaitCount))
				DBPoolWaitDuration.Set(stats.WaitDuration.Seconds())
			}
		}
	}()
}
