package stats

import (
	"context"
	"database/sql"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"
)

func query(t *testing.T, db *sql.DB) {
	for i := 0; i < 10; i++ {
		rows, err := db.Query("SELECT * FROM accounts")
		if err != nil {
			t.Fatal(err)
		}
		rows.Close()
	}
}

func TestGetProjectsHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/metrics", nil)
	if err != nil {
		t.Fatal(err)
	}
	db, err := sql.Open("postgres", "postgres://pinger:pinger@localhost:5432/pinger")
	if err != nil {
		t.Fatal(err)
	}

	collector := NewSQLStats("db_name", db)

	prometheus.MustRegister(collector)
	query(t, db)
	rr := httptest.NewRecorder()
	handler := promhttp.Handler()

	// Populate the request's context with our test data.
	ctx := req.Context()
	ctx = context.WithValue(ctx, "testrequest", "testing")

	req = req.WithContext(ctx)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEqual(t, -1, strings.Index(string(body), "db_stats"))

}
