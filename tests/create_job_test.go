package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestCreateJob(t *testing.T) {
	server := SetupTestApp(t)
	defer Teardown(t, server)

	job := map[string]interface{}{
		"title":       "Backend Developer",
		"description": "Build scalable APIs",
		"location":    "Lagos, Nigeria",
		"salary":      120000.0,
		"duties":      []string{"Write code", "Review PRs"},
		"url":         "https://example.com/job/1",
	}

	body, _ := json.Marshal(job)
	resp, err := http.Post(server.URL+"/jobs", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}

	// Parse response
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	// Verify response
	if result["message"] != "job created" {
		t.Errorf("Expected 'job created', got '%v'", result["message"])
	}

	jobData := result["job"].(map[string]interface{})
	if jobData["title"] != job["title"] {
		t.Errorf("Expected title '%s', got '%s'", job["title"], jobData["title"])
	}

	if jobData["id"] == nil {
		t.Error("Job should have an ID")
	}
}

func TestCreateJob_MissingTitle(t *testing.T) {
	server := SetupTestApp(t)
	defer Teardown(t, server)

	job := map[string]interface{}{
		"description": "Build APIs",
		"location":    "USA",
		"salary":      120000.0,
		"duties":      []string{"Write code"},
		"url":         "https://example.com/job/1",
	}

	body, _ := json.Marshal(job)
	resp, err := http.Post(server.URL+"/jobs", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}
