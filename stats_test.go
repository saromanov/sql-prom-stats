package stats

import (
	"database/sql"
	"fmt"
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

func idleConnectionQuery(t *testing.T, db *sql.DB) {
	for i := 0; i < 100; i++ {
		rows, err := db.Query("SELECT * FROM accounts")
		if err != nil {
			t.Fatal(err)
		}
		rows.Close()
	}
}

func TestGetProjectsHandler(t *testing.T) {
	db, err := sql.Open("postgres", "postgres://pinger:pinger@localhost:5432/pinger")
	if err != nil {
		t.Fatal(err)
	}

	collector := NewSQLStats("db_name", db)

	prometheus.MustRegister(collector)
	srv := httptest.NewServer(promhttp.Handler())
	defer srv.Close()
	idleConnectionQuery(t, db)
	rr, err := http.Get(fmt.Sprintf("%s/metrics", srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	if status := rr.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, string(body), "AA")
	assert.NotEqual(t, -1, strings.Index(string(body), "idle"))

}
