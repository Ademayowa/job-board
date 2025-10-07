package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

// TestGetAllJobs tests for fetching all jobs
func TestGetAllJobs(t *testing.T) {
	server := SetupTestApp(t)
	defer Teardown(t, server)

	// Create test jobs
	job1 := map[string]interface{}{
		"title":       "Backend Developer",
		"description": "Build APIs",
		"location":    "Australia",
		"salary":      120000.0,
		"duties":      []string{"Code", "Review"},
		"url":         "http://example.com/1",
	}

	job2 := map[string]interface{}{
		"title":       "Frontend Developer",
		"description": "Build UIs",
		"location":    "Remote",
		"salary":      100000.0,
		"duties":      []string{"Design", "Code"},
		"url":         "http://example.com/2",
	}

	// Create jobs via API
	body1, _ := json.Marshal(job1)
	http.Post(server.URL+"/jobs", "application/json", bytes.NewBuffer(body1))

	body2, _ := json.Marshal(job2)
	http.Post(server.URL+"/jobs", "application/json", bytes.NewBuffer(body2))

	// Fetch all jobs
	resp, err := http.Get(server.URL + "/jobs")
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Parse response
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	// Verify data
	data := result["data"].([]interface{})
	if len(data) != 2 {
		t.Errorf("Expected 2 jobs, got %d", len(data))
	}

	// Verify metadata
	metadata := result["metadata"].(map[string]interface{})
	if metadata["total"] != float64(2) {
		t.Errorf("Expected total 2, got %v", metadata["total"])
	}
}
