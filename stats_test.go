package stats

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
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

func getStatNumber(s, prefix string) int {
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "database_sql_stats_golang_idle") {
			lastChar := line[len(line)-1 : len(line)]
			res, _ := strconv.ParseInt(lastChar, 10, 32)
			return int(res)
		}
	}

	return 0
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
	data := string(body)
	idle := getStatNumber(data, "database_sql_stats_golang_idle")
	assert.NotEqual(t, -1, strings.Index(data, "database_sql_stats_golang_max_idle"))
	assert.Equal(t, 1, idle)

}
