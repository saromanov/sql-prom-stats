package stats

import (
	"sync"
	"database/sql"
	"time"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	initPromMetricsOnce = &sync.Once{}
	promMetric          = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "db_stats",
			Help: "SQL database stats",
		},
		[]string{"db_name", "metric"},
	)
)

// collector defines collecting of metrics
// to prometheus
type collector struct{ 
	dbName string
	getter StatsGetter
}

func newPrometheus(dbName string) *collector {
	initPromMetricsOnce.Do(func() { pr.MustRegister(promMetric) })
	return &collector{dbName: dbName}
}


// StatsGetter interface for getting stats
// from sql.DB
type StatsGetter interface {
	Stats() sql.DBStats
}

func (p *collector) Collect(stats sql.DBStats) {
	
}


func StartCollectPrometheusMetrics(db StatsGetter, interval time.Duration, dbName string) CollectorStopper {
	return StartCollect(db, interval, newPrometheus(dbName))
}

