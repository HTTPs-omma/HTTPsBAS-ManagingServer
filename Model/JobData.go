package Model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type JobData struct {
	Id          int       `json:"id"`
	ProcedureID string    `json:"procedureID"`
	AgentUUID   string    `json:"agentUUID"`
	MessageUUID string    `json:"messageUUID"`
	Upload      string    `json:"upload"` // Update를 Upload로 변경
	Action      string    `json:"action"`
	Files       []string  `json:"files"`
	CreateAt    time.Time `json:"createAt"`
}

type JobDB struct {
	dbName string
}

func NewJobDB() (*JobDB, error) {
	jd := &JobDB{dbName: "jobs"}
	err := jd.createTableIfNotExists()
	if err != nil {
		return nil, err
	}
	return jd, nil
}

func (jd *JobDB) createTableIfNotExists() error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS jobs (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
		ProcedureID TEXT,
		AgentUUID TEXT,
		MessageUUID TEXT,
		Upload TEXT,   -- Update 필드를 Upload로 변경
		Action TEXT,
		Files TEXT,
		CreateAt DATETIME
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("create table failed: %w", err)
	}
	return nil
}

func (jd *JobDB) InsertJobData(jobData *JobData) error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	jobData.CreateAt = time.Now()

	// Upload 배열을 JSON 문자열로 직렬화
	uploadJSON, err := json.Marshal(jobData.Upload)
	if err != nil {
		return fmt.Errorf("failed to marshal upload: %w", err)
	}

	// Files 배열을 JSON 문자열로 직렬화
	filesJSON, err := json.Marshal(jobData.Files)
	if err != nil {
		return fmt.Errorf("failed to marshal files: %w", err)
	}

	insertSQL := `INSERT INTO jobs (ProcedureID, AgentUUID, MessageUUID, Upload, Action, Files, CreateAt) 
	              VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err = db.Exec(insertSQL, jobData.ProcedureID, jobData.AgentUUID, jobData.MessageUUID, string(uploadJSON), jobData.Action, string(filesJSON), jobData.CreateAt)
	if err != nil {
		return fmt.Errorf("insert job failed: %w", err)
	}
	return nil
}

func (jd *JobDB) SelectJobDataByAgentUUID(agentUUID string) (*JobData, error, bool) {
	db, err := getDBPtr()
	if err != nil {
		return nil, err, false
	}
	defer db.Close()

	selectSQL := `SELECT id, ProcedureID, AgentUUID, MessageUUID, Upload, Action, Files, CreateAt FROM jobs WHERE AgentUUID = ?`
	rows, err := db.Query(selectSQL, agentUUID)
	if err != nil {
		return nil, err, false
	}
	defer rows.Close()

	if rows.Next() == false {
		return &JobData{}, nil, false
	}

	var job *JobData = &JobData{}
	var uploadJSON string
	var filesJSON string

	err = rows.Scan(&job.Id, &job.ProcedureID, &job.AgentUUID, &job.MessageUUID, &uploadJSON, &job.Action, &filesJSON, &job.CreateAt)
	if err != nil {
		return nil, err, true
	}

	// Upload JSON 문자열을 []string 배열로 역직렬화
	err = json.Unmarshal([]byte(uploadJSON), &job.Upload)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal upload: %w", err), true
	}

	// Files JSON 문자열을 []string 배열로 역직렬화
	err = json.Unmarshal([]byte(filesJSON), &job.Files)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal files: %w", err), true
	}

	return job, nil, true
}

func (jd *JobDB) SelectAllJobData() ([]JobData, error) {
	db, err := getDBPtr()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	selectSQL := `SELECT id, ProcedureID, AgentUUID, MessageUUID, Upload, Action, Files, CreateAt FROM jobs`
	rows, err := db.Query(selectSQL)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	if rows.Next() == false {
		return []JobData{}, nil
	}

	jobs := []JobData{}

	for {
		var job JobData
		var uploadJSON string
		var filesJSON string

		err = rows.Scan(&job.Id, &job.ProcedureID, &job.AgentUUID, &job.MessageUUID, &uploadJSON, &job.Action, &filesJSON, &job.CreateAt)
		if err != nil {
			return nil, err
		}

		// Upload JSON 문자열을 []string 배열로 역직렬화
		err = json.Unmarshal([]byte(uploadJSON), &job.Upload)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal upload: %w", err)
		}

		// Files JSON 문자열을 []string 배열로 역직렬화
		err = json.Unmarshal([]byte(filesJSON), &job.Files)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal files: %w", err)
		}

		jobs = append(jobs, job)

		if !rows.Next() {
			break
		}
	}

	return jobs, nil
}

func (jd *JobDB) PopbyAgentUUID(agentUUID string) (*JobData, error, bool) {
	db, err := getDBPtr()
	if err != nil {
		return nil, err, false
	}
	defer db.Close()

	selectSQL := `
		SELECT id, ProcedureID, AgentUUID, MessageUUID, Upload, Action, Files, CreateAt 
		FROM jobs 
		WHERE AgentUUID = ? 
		ORDER BY CreateAt ASC
		LIMIT 1
	`
	row := db.QueryRow(selectSQL, agentUUID)

	var job JobData
	var uploadJSON string
	var filesJSON string

	err = row.Scan(&job.Id, &job.ProcedureID, &job.AgentUUID, &job.MessageUUID, &uploadJSON, &job.Action, &filesJSON, &job.CreateAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return &JobData{}, nil, false
		}
		return nil, err, false
	}

	// Upload JSON 문자열을 []string 배열로 역직렬화
	err = json.Unmarshal([]byte(uploadJSON), &job.Upload)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal upload: %w", err), false
	}

	// Files JSON 문자열을 []string 배열로 역직렬화
	err = json.Unmarshal([]byte(filesJSON), &job.Files)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal files: %w", err), false
	}

	err = jd.DeleteJobDataById(job.Id)
	if err != nil {
		return nil, err, false
	}

	return &job, nil, true
}

// DeleteJobDataById: ID 기반으로 JobData 삭제
func (jd *JobDB) DeleteJobDataById(id int) error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	deleteSQL := `DELETE FROM jobs WHERE id = ?`
	_, err = db.Exec(deleteSQL, id)
	if err != nil {
		return fmt.Errorf("delete job failed: %w", err)
	}
	return nil
}

func (jd *JobDB) DeleteAllJobData() error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	deleteSQL := "DELETE FROM jobs"
	_, err = db.Exec(deleteSQL)
	if err != nil {
		return fmt.Errorf("delete job failed: %w", err)
	}
	return nil
}
