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
	inUseDesc       *prometheus.Desc
	idleDesc        *prometheus.Desc
	maxIdleDesc     *prometheus.Desc
	maxLifetimeDesc *prometheus.Desc
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
		maxLifetimeDesc: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "max_idle"), "The total number of connections closed due to SetConnMaxLifetime.", []string{"db"}, nil),
		inUseDesc:       prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "max_idle"), "The number of connections currently in use.", []string{"db"}, nil),
		idleDesc:        prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "max_idle"), "The number of idle connections.", []string{"db"}, nil),
		maxOpenDesc:     prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "max_idle"), "Maximum number of open connections to the database.", []string{"db"}, nil),
		openDesc:        prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "max_idle"), "The number of established connections both in use and idle.", []string{"db"}, nil),
		waitedForDesc:   prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "max_idle"), "The total number of connections waited for.", []string{"db"}, nil),
	}
}

func (p *collector) Collect(ch chan<- prometheus.Metric) {
	stats := p.getter.Stats()

	ch <- prometheus.MustNewConstMetric(
		p.maxOpenDesc,
		prometheus.GaugeValue,
		float64(stats.MaxOpenConnections),
		p.dbName,
	)
	ch <- prometheus.MustNewConstMetric(
		p.openDesc,
		prometheus.GaugeValue,
		float64(stats.OpenConnections),
		p.dbName,
	)
	ch <- prometheus.MustNewConstMetric(
		p.inUseDesc,
		prometheus.GaugeValue,
		float64(stats.InUse),
		p.dbName,
	)
	ch <- prometheus.MustNewConstMetric(
		p.idleDesc,
		prometheus.GaugeValue,
		float64(stats.Idle),
		p.dbName,
	)
	ch <- prometheus.MustNewConstMetric(
		p.waitedForDesc,
		prometheus.CounterValue,
		float64(stats.WaitCount),
		p.dbName,
	)
	ch <- prometheus.MustNewConstMetric(
		p.maxIdleDesc,
		prometheus.CounterValue,
		float64(stats.MaxIdleClosed),
		p.dbName,
	)
	ch <- prometheus.MustNewConstMetric(
		p.maxLifetimeDesc,
		prometheus.CounterValue,
		float64(stats.MaxLifetimeClosed),
		p.dbName,
	)
}
