package Model

import (
	"fmt"
	"time"
)

type JobData struct {
	ProcedureID     string    `json:"procedureID"`
	AgentUUID       string    `json:"agentUUID"`
	InstructionUUID string    `json:"instructionUUID"`
	CreateAt        time.Time `json:"createAt"`
}

/**
Chatgpt 이용해서 생성
*/

type JobDB struct {
	dbName string
}

// NewJobDB: SQLite DB를 초기화하고 테이블을 생성하는 함수
func NewJobDB() (*JobDB, error) {
	jd := &JobDB{dbName: "jobs"}
	err := jd.createTableIfNotExists()
	if err != nil {
		return nil, err
	}
	return jd, nil
}

// createTableIfNotExists: jobs 테이블이 없으면 생성
func (jd *JobDB) createTableIfNotExists() error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS jobs (
		ProcedureID TEXT,
		AgentUUID TEXT,
		InstructionUUID TEXT,
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
	insertSQL := `INSERT INTO jobs (ProcedureID, AgentUUID, InstructionUUID, CreateAt) 
	              VALUES (?, ?, ?, ?)`
	_, err = db.Exec(insertSQL, jobData.ProcedureID, jobData.AgentUUID, jobData.InstructionUUID, jobData.CreateAt)
	fmt.Println(jobData)
	if err != nil {
		return fmt.Errorf("insert job failed: %w", err)
	}
	return nil
}

// GetJobDataByAgentUUID: AgentUUID 기반으로 JobData 조회
func (jd *JobDB) SelectJobDataByAgentUUID(agentUUID string) (*JobData, error, bool) {
	db, err := getDBPtr()
	if err != nil {
		return nil, err, false
	}
	defer db.Close()

	selectSQL := `SELECT ProcedureID, AgentUUID, InstructionUUID, CreateAt FROM jobs WHERE AgentUUID = ?`
	rows, err := db.Query(selectSQL, agentUUID)
	if err != nil {
		return nil, err, false
	}
	defer rows.Close()
	if rows.Next() == false {
		// 결과가 없다면,
		return &JobData{}, nil, false
	}

	var job *JobData = &JobData{}
	err = rows.Scan(&job.ProcedureID, &job.AgentUUID, &job.InstructionUUID, &job.CreateAt) // 첫 행에만 적용
	if err != nil {
		return nil, err, true
	}

	return job, nil, true
}

func (jd *JobDB) SelectAllJobData() ([]JobData, error) {

	db, err := getDBPtr()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	selectSQL := `SELECT ProcedureID, AgentUUID, InstructionUUID, CreateAt FROM jobs`

	rows, err := db.Query(selectSQL)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	if rows.Next() == false {
		return []JobData{}, nil
	}

	jobs := []JobData{}
	for rows.Next() == true {
		var job *JobData = &JobData{}
		err = rows.Scan(&job.ProcedureID, &job.AgentUUID, &job.InstructionUUID, &job.CreateAt) // 첫 행에만 적용
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, *job)
	}

	return jobs, nil
}

// DeleteJobDataByInstructionUUID: InstructionUUID 기반으로 JobData 삭제
func (jd *JobDB) DeleteJobDataByInstructionUUID(instructionUUID string) error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	deleteSQL := `DELETE FROM jobs WHERE InstructionUUID = ?`
	_, err = db.Exec(deleteSQL, instructionUUID)
	if err != nil {
		return fmt.Errorf("delete job failed: %w", err)
	}
	return nil
}

// DeleteJobDataByInstructionUUID: InstructionUUID 기반으로 JobData 삭제
func (jd *JobDB) DeleteAllJobData() error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	deleteSQL := `DELETE FROM jobs`
	_, err = db.Exec(deleteSQL)
	if err != nil {
		return fmt.Errorf("delete job failed: %w", err)
	}
	return nil
}
