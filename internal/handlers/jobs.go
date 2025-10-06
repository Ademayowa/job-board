package handlers

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	"github.com/Ademayowa/job-board/internal/config"
	"github.com/Ademayowa/job-board/internal/models"

	"github.com/gin-gonic/gin"
)

// Create a job
func createJob(context *gin.Context) {
	var job models.Job

	err := context.ShouldBindJSON(&job)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "could not parse job data"})
		return
	}

	job.Save()
	context.JSON(http.StatusCreated, gin.H{"message": "job created", "job": job})
}

// Fetch all jobs
func getJobs(context *gin.Context) {
	// Extract job query parameter from the URL
	filterTitle := context.Query("query")

	// Extract pagination parameters with defaults
	page, err := strconv.Atoi(context.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1 // Default to page 1 if invalid
	}

	limit, err := strconv.Atoi(context.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 6 // Default to 10 items per page if invalid
	}

	// Get all jobs with filters and pagination
	jobs, total, err := models.GetAllJobs(filterTitle, page, limit)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch jobs"})
		return
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	// Return jobs with the metadata(all jobs in the database & pagination)
	context.JSON(http.StatusOK, gin.H{
		"data": jobs,
		"metadata": gin.H{
			"current_page": page,
			"per_page":     limit,
			"total":        total,
			"total_pages":  totalPages,
		},
	})
}

// Fetch a single job
func getJob(context *gin.Context) {
	jobId := context.Param("id")

	job, err := models.GetJobByID(jobId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch job"})
		return
	}

	context.JSON(http.StatusOK, job)
}

// Delete a job
func deleteJob(context *gin.Context) {
	jobId := context.Param("id")

	job, err := models.GetJobByID(jobId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch job"})
		return
	}

	err = job.Delete()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete job"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "job deleted successfully"})
}

// Update a job
func updateJob(context *gin.Context) {
	// Extract job ID from the URL
	jobId := context.Param("id")

	// Parse the request body to get the updated job data
	var updatedJob models.Job
	if err := context.ShouldBindJSON(&updatedJob); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body" + err.Error()})
		return
	}

	// Convert Duties field to JSON for database storage
	dutiesJSON, err := json.Marshal(updatedJob.Duties)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "error processing duties field"})
		return
	}

	// Update job in the database
	err = models.UpdateJobByID(jobId, updatedJob, string(dutiesJSON))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "could not update job"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "job updated successfully"})
}

// Get jobs sorted by most recent
func GetRecentJobs(context *gin.Context) {
	limitParam := context.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitParam)

	jobs, err := models.GetJobsSortedByRecent(limit)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch recent jobs"})
		return
	}

	context.JSON(http.StatusOK, jobs)
}

// Get jobs sorted by highest salary
func GetHighestSalaryJobs(context *gin.Context) {
	limitParam := context.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitParam)

	jobs, err := models.GetJobsSortedBySalary(limit)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch highest salary jobs"})
		return
	}

	context.JSON(http.StatusOK, jobs)
}

// ShareJobLink returns a shareable link for a job post
func ShareJobLink(context *gin.Context) {
	jobId := context.Param("id")

	baseURL := context.Request.Host
	if baseURL == "" {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "could not determine host URL"})
		return
	}
	// Check if the app runs on http or https
	scheme := "http"
	if context.Request.TLS != nil {
		scheme = "https"
	}
	// Generate a shareable link to the job details page
	shareableLink := scheme + "://" + baseURL + config.JobDetailsPage + "/" + jobId

	context.JSON(http.StatusOK, shareableLink)
}
