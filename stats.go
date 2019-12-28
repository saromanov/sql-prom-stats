package stats

import (
	"database/sql"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "database_sql_stats"
	subsystem = "golang"
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

// StatsGetter interface for getting stats
// from sql.DB
type StatsGetter interface {
	Stats() sql.DBStats
}

// collector defines collecting of metrics
// to prometheus
type collector struct {
	dbName          string
	getter          StatsGetter
	maxIdleDesc     *prometheus.Desc
	maxLifetimeDesc *prometheus.Desc
	inUseDesc       *prometheus.Desc
	idleDesc        *prometheus.Desc
	maxOpenDesc     *prometheus.Desc
	openDesc        *prometheus.Desc
	waitedForDesc   *prometheus.Desc
}

func newPrometheus(dbName string, getter sql.DBStats) *collector {
	initPromMetricsOnce.Do(func() { prometheus.MustRegister(promMetric) })
	return &collector{
		dbName:          dbName,
		getter:          getter,
		maxIdleDesc:     prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "max_idle"), "The total number of connections closed due to SetMaxIdleConns", []string{"db"}, nil),
		maxLifetimeDesc: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "max_idle"), "The total number of connections closed due to SetMaxIdleConns", []string{"db"}, nil),
		inUseDesc:       prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "max_idle"), "The total number of connections closed due to SetMaxIdleConns", []string{"db"}, nil),
		idleDesc:        prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "max_idle"), "The total number of connections closed due to SetMaxIdleConns", []string{"db"}, nil),
		maxOpenDesc:     prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "max_idle"), "The total number of connections closed due to SetMaxIdleConns", []string{"db"}, nil),
		openDesc:        prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "max_idle"), "The total number of connections closed due to SetMaxIdleConns", []string{"db"}, nil),
		waitedForDesc:   prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "max_idle"), "The total number of connections closed due to SetMaxIdleConns", []string{"db"}, nil),
	}
}

func (p *collector) Collect(stats sql.DBStats) {

}

func StartCollect() {

}
