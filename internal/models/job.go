package models

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/Ademayowa/job-board/internal/database/db"

	"github.com/google/uuid"
)

// DateFormat is the standard date format used throughout the application
const DateFormat = time.RFC3339

type Job struct {
	ID          string   `json:"id"`
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Location    string   `json:"location" binding:"required"`
	Salary      float64  `json:"salary" binding:"required"`
	Duties      []string `json:"duties" binding:"required"`
	Url         string   `json:"url"`
	CreatedAt   string   `json:"created_at"`
	Expired     bool     `json:"expired"`
}

// IsExpired checks if a job is expired
func (job *Job) IsExpired() bool {
	return job.DaysToExpiration() <= 0
}

// DaysToExpiration returns the number of days until job expires
// Positive: days remaining, Zero: expires today, Negative: days since expiration
func (job *Job) DaysToExpiration() int {
	createdAt, err := time.Parse(DateFormat, job.CreatedAt)
	if err != nil {
		return 0
	}

	expirationDate := createdAt.AddDate(0, 0, 14) // 14 days after creation

	// Calculate days remaining
	now := time.Now()
	duration := expirationDate.Sub(now)
	return int(duration.Hours() / 24)
}

// Save job into the database
func (job *Job) Save() error {
	job.ID = uuid.New().String()

	dutiesJSON, err := json.Marshal(job.Duties)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO jobs(id, title, description, location, salary, duties, url, created_at)
		VALUES(?, ?, ?, ?, ?, ?, ?, ?)
	`

	sqlStmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer sqlStmt.Close()

	job.CreatedAt = time.Now().Format(DateFormat)

	_, err = sqlStmt.Exec(
		job.ID,
		job.Title,
		job.Description,
		job.Location,
		job.Salary,
		string(dutiesJSON),
		job.Url,
		job.CreatedAt,
	)

	return err
}

// Get all jobs (with optional filtering by job title)
func GetAllJobs(filterTitle string, page, limit int) ([]Job, int, error) {
	query := "SELECT * FROM jobs WHERE 1=1"
	args := []interface{}{}

	// Filter jobs by the title
	if strings.TrimSpace(filterTitle) != "" {
		query += " AND LOWER(title) LIKE ?"
		args = append(args, "%"+strings.ToLower(filterTitle)+"%")
	}

	// Count total jobs that matches the filter from the database
	countQuery := "SELECT COUNT(*) FROM (" + query + ") AS count_query"

	var total int
	err := db.DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Add pagination
	offset := (page - 1) * limit
	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	// Fetch paginated jobs
	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var jobs []Job

	for rows.Next() {
		var job Job
		var dutiesJSON string

		err := rows.Scan(
			&job.ID,
			&job.Title,
			&job.Description,
			&job.Location,
			&job.Salary,
			&dutiesJSON,
			&job.Url,
			&job.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		// Convert Duties field to []string
		if err := json.Unmarshal([]byte(dutiesJSON), &job.Duties); err != nil {
			return nil, 0, err
		}

		// Check if job is expired
		job.Expired = job.IsExpired()

		jobs = append(jobs, job)
	}

	return jobs, total, nil
}

// Get a job by ID
func GetJobByID(id string) (Job, error) {
	var job Job
	var dutiesJSON string

	query := "SELECT * FROM jobs WHERE id =?"
	row := db.DB.QueryRow(query, id)

	err := row.Scan(
		&job.ID,
		&job.Title,
		&job.Description,
		&job.Location,
		&job.Salary, &dutiesJSON,
		&job.Url,
		&job.CreatedAt,
	)
	if err != nil {
		return job, err
	}

	// Convert Duties field from JSON to []string
	err = json.Unmarshal([]byte(dutiesJSON), &job.Duties)
	if err != nil {
		return job, err
	}

	// Check if job is expired
	job.Expired = job.IsExpired()

	return job, nil
}

// Delete a job
func (job Job) Delete() error {
	query := "DELETE FROM jobs WHERE id = ?"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(job.ID)

	return err
}

// Update a job by ID
func UpdateJobByID(id string, updatedJob Job, dutiesJSON string) error {
	query := `
		UPDATE jobs
		SET title = ?, description = ?, location = ?, salary = ?, duties = ?, url = ?
		WHERE id = ?
	`
	_, err := db.DB.Exec(query,
		updatedJob.Title,
		updatedJob.Description,
		updatedJob.Location,
		updatedJob.Salary,
		dutiesJSON,
		updatedJob.Url,
		id,
	)

	return err
}

// Get jobs sorted by most recent
func GetJobsSortedByRecent(limit int) ([]Job, error) {
	query := "SELECT * FROM jobs ORDER BY created_at DESC LIMIT ?"

	rows, err := db.DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []Job

	for rows.Next() {
		var job Job
		var dutiesJSON string

		err := rows.Scan(
			&job.ID,
			&job.Title,
			&job.Description,
			&job.Location,
			&job.Salary,
			&dutiesJSON,
			&job.Url,
			&job.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(dutiesJSON), &job.Duties)
		if err != nil {
			return nil, err
		}

		// Check if job is expired
		job.Expired = job.IsExpired()

		jobs = append(jobs, job)
	}

	return jobs, nil
}

// Get jobs sorted by highest salary
func GetJobsSortedBySalary(limit int) ([]Job, error) {
	query := "SELECT * FROM jobs ORDER BY salary DESC LIMIT ?"

	rows, err := db.DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []Job

	for rows.Next() {
		var job Job
		var dutiesJSON string

		err := rows.Scan(
			&job.ID,
			&job.Title,
			&job.Description,
			&job.Location,
			&job.Salary,
			&dutiesJSON,
			&job.Url,
			&job.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(dutiesJSON), &job.Duties)
		if err != nil {
			return nil, err
		}

		// Check if job is expired
		job.Expired = job.IsExpired()

		jobs = append(jobs, job)
	}

	return jobs, nil
}
