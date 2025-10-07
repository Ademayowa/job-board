package tests

import (
	"database/sql"
	"net/http/httptest"
	"testing"

	db "github.com/Ademayowa/job-board/internal/database"
	routes "github.com/Ademayowa/job-board/internal/handlers"
	"github.com/gin-gonic/gin"
)

// SetupTestApp sets up the test environment
func SetupTestApp(t *testing.T) *httptest.Server {
	gin.SetMode(gin.TestMode)

	// Setup in-memory test database
	var err error
	db.DB, err = sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Create jobs table
	createTable := `
	CREATE TABLE IF NOT EXISTS jobs (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		salary FLOAT NOT NULL,
		duties TEXT NOT NULL,
		url TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	if _, err = db.DB.Exec(createTable); err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// Setup router
	router := gin.New()
	routes.RegisterRoutes(router)

	return httptest.NewServer(router)
}

func Teardown(t *testing.T, server *httptest.Server) {
	server.Close()
	if db.DB != nil {
		db.DB.Close()
	}
}
