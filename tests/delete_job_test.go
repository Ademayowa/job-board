package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

// TestDeleteJob tests for deleting a job
func TestDeleteJob(t *testing.T) {
	server := SetupTestApp(t)
	defer Teardown(t, server)

	// Create a job
	job := map[string]interface{}{
		"title":       "Backend Developer",
		"description": "Build APIs",
		"location":    "Lagos",
		"salary":      120000.0,
		"duties":      []string{"Write code"},
		"url":         "http://example.com/job/1",
	}

	body, _ := json.Marshal(job)
	createResp, err := http.Post(server.URL+"/jobs", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create job: %v", err)
	}
	defer createResp.Body.Close()

	// Get the created job ID
	var createResult map[string]interface{}
	json.NewDecoder(createResp.Body).Decode(&createResult)
	jobData := createResult["job"].(map[string]interface{})
	jobID := jobData["id"].(string)

	// Delete the job
	req, _ := http.NewRequest("DELETE", server.URL+"/jobs/"+jobID, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Verify response message
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if result["message"] != "job deleted successfully" {
		t.Errorf("Expected 'job deleted successfully', got '%v'", result["message"])
	}

	// Verify job is deleted (fetching should fail)
	getResp, _ := http.Get(server.URL + "/jobs/" + jobID)
	if err != nil {
		t.Fatalf("Failed to fetch job: %v", err)
	}
	defer getResp.Body.Close()

	if getResp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status 500 for deleted job, got %d", getResp.StatusCode)
	}
}

func TestDeleteJob_NotFound(t *testing.T) {
	server := SetupTestApp(t)
	defer Teardown(t, server)

	// Try to delete non-existent job
	req, _ := http.NewRequest("DELETE", server.URL+"/jobs/invalid-id", nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", resp.StatusCode)
	}
}
