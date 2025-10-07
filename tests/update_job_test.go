// tests/update_job_test.go
package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

// TestUpdateJob tests for updating a job
func TestUpdateJob(t *testing.T) {
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

	// Update the job
	updatedJob := map[string]interface{}{
		"title":       "DevOps Engineer",
		"description": "Build CI/CD pipelines",
		"location":    "Remote",
		"salary":      150000.0,
		"duties":      []string{"Implement CI/CD pipelines", "Monitor infrastructure"},
		"url":         "http://example.com/job/updated",
	}

	updateBody, _ := json.Marshal(updatedJob)
	req, _ := http.NewRequest("PUT", server.URL+"/jobs/"+jobID, bytes.NewBuffer(updateBody))
	req.Header.Set("Content-Type", "application/json")

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

	if result["message"] != "job updated successfully" {
		t.Errorf("Expected 'job updated successfully', got '%v'", result["message"])
	}

	// Fetch the updated job to verify
	getResp, err := http.Get(server.URL + "/jobs/" + jobID)
	if err != nil {
		t.Fatalf("Failed to fetch job: %v", err)
	}
	defer getResp.Body.Close()

	var fetchedJob map[string]interface{}
	json.NewDecoder(getResp.Body).Decode(&fetchedJob)

	if fetchedJob["title"] != updatedJob["title"] {
		t.Errorf("Expected title '%s', got '%s'", updatedJob["title"], fetchedJob["title"])
	}
}
