package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

// TestGetSingleJob tests for fetching a single job
func TestGetSingleJob(t *testing.T) {
	server := SetupTestApp(t)
	defer Teardown(t, server)

	// Create a job
	job := map[string]interface{}{
		"title":       "Backend Developer",
		"description": "Build scalable APIs",
		"location":    "Lagos",
		"salary":      120000.0,
		"duties":      []string{"Write code", "Review PRs"},
		"url":         "http://example.com/job/1",
	}

	body, _ := json.Marshal(job)
	createResp, _ := http.Post(server.URL+"/jobs", "application/json", bytes.NewBuffer(body))
	defer createResp.Body.Close()

	// Get the created job ID
	var createResult map[string]interface{}
	json.NewDecoder(createResp.Body).Decode(&createResult)
	jobData := createResult["job"].(map[string]interface{})
	jobID := jobData["id"].(string)

	// Fetch the single job
	resp, err := http.Get(server.URL + "/jobs/" + jobID)
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

	// Verify job data
	if result["id"] != jobID {
		t.Errorf("Expected job ID %s, got %s", jobID, result["id"])
	}

	if result["title"] != job["title"] {
		t.Errorf("Expected title '%s', got '%s'", job["title"], result["title"])
	}

	if result["location"] != job["location"] {
		t.Errorf("Expected location '%s', got '%s'", job["location"], result["location"])
	}
}

func TestGetSingleJob_NotFound(t *testing.T) {
	server := SetupTestApp(t)
	defer Teardown(t, server)

	// Try to fetch non-existent job
	resp, err := http.Get(server.URL + "/jobs/invalid-id")
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", resp.StatusCode)
	}
}
