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

// NewSQLStats provides initialization of collecting of metrics
func NewSQLStats(dbName string, getter StatsGetter) *collector {
	initPromMetricsOnce.Do(func() { prometheus.MustRegister(promMetric) })
	return &collector{
		dbName:          dbName,
		getter:          getter,
		maxIdleDesc:     prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "max_idle"), "The total number of connections closed due to SetMaxIdleConns", []string{"db_stat"}, nil),
		maxLifetimeDesc: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "max_lifetime"), "The total number of connections closed due to SetConnMaxLifetime.", []string{"db_stat"}, nil),
		inUseDesc:       prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "in_use"), "The number of connections currently in use.", []string{"db_stat"}, nil),
		idleDesc:        prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "idle"), "The number of idle connections.", []string{"db_stat"}, nil),
		maxOpenDesc:     prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "max_open"), "Maximum number of open connections to the database.", []string{"db_stat"}, nil),
		openDesc:        prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "open_connections"), "The number of established connections both in use and idle.", []string{"db_stat"}, nil),
		waitedForDesc:   prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "waited_for"), "The total number of connections waited for.", []string{"db_stat"}, nil),
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

func (p *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- p.maxOpenDesc
	ch <- p.openDesc
	ch <- p.inUseDesc
	ch <- p.idleDesc
	ch <- p.waitedForDesc
	ch <- p.maxIdleDesc
	ch <- p.maxLifetimeDesc
}
