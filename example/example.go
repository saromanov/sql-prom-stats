package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	stats "github.com/saromanov/sql-prom-stats"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func query(db *sql.DB) {
	for i := 0; i < 100; i++ {
		rows, err := db.Query("SELECT * FROM accounts")
		if err != nil {
			fmt.Println(err)
		}
		rows.Close()
	}
}

func run() error {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/postgres")
	if err != nil {
		return err
	}

	go query(db)
	collector := stats.NewSQLStats("db_name", db)
	prometheus.MustRegister(collector)
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(":8080", nil)
}
