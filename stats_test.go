package stats

import (
	"testing"
	"net/http/httptest"
	"net/http"

	"github.com/stretchr/testify/assert"
)

func TestGetProjectsHandler(t *testing.T) {
    req, err := http.NewRequest("GET", "/metrics", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(GetUsersHandler)

    // Populate the request's context with our test data.
    ctx := req.Context()
	ctx = context.WithValue(ctx, "testrequest", "testing")
	
    req = req.WithContext(ctx)
    handler.ServeHTTP(rr, req)

    // Check the status code is what we expect.
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }
}