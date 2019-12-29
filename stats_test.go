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

func leakConnectionsQuery(t *testing.T, db *sql.DB) {
	for i := 0; i < 10; i++ {
		_, err := db.Query("SELECT * FROM accounts")
		if err != nil {
			t.Fatal(err)
		}
	}
}

func getStatNumber(s, prefix string) int {
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, prefix) {
			dataSlice := strings.Split(line, " ")
			lastChar := dataSlice[len(dataSlice)-1]
			res, _ := strconv.ParseInt(lastChar, 10, 32)
			return int(res)
		}
	}

	return 0
}

func makeRequest(url string) (string, error) {
	rr, err := http.Get(fmt.Sprintf("%s/metrics", url))
	if err != nil {
		return "", err
	}
	if status := rr.StatusCode; status != http.StatusOK {
		return "", fmt.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
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
	data, err := makeRequest(srv.URL)
	assert.NoError(t, err)
	idle := getStatNumber(data, "database_sql_stats_golang_idle")
	assert.NotEqual(t, -1, strings.Index(data, "database_sql_stats_golang_max_idle"))
	assert.Equal(t, 1, idle)

	db.Close()
	data, err = makeRequest(srv.URL)
	assert.NoError(t, err)
	assert.Equal(t, 0, getStatNumber(data, "database_sql_stats_golang_idle"))

}

func TestGetProjectsHandlerLeakConnections(t *testing.T) {
	db, err := sql.Open("postgres", "postgres://pinger:pinger@localhost:5432/pinger")
	if err != nil {
		t.Fatal(err)
	}

	collector := NewSQLStats("db_name", db)

	prometheus.MustRegister(collector)
	srv := httptest.NewServer(promhttp.Handler())
	defer srv.Close()
	leakConnectionsQuery(t, db)
	data, err := makeRequest(srv.URL)
	assert.NoError(t, err)
	inUse := getStatNumber(data, "database_sql_stats_golang_in_use")
	assert.NotEqual(t, -1, strings.Index(data, "database_sql_stats_golang_in_use"))
	assert.Equal(t, 10, inUse)
}
