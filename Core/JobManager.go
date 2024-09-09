package Core

import (
	"fmt"
	"github.com/your/repo/Model"
)

type JobManager struct {
	jobDB *Model.JobDB
}

// NewJobManager: JobDB와 연결하고 JobManager 생성
func NewJobManager() (*JobManager, error) {
	jobDB, err := Model.NewJobDB()
	if err != nil {
		return nil, err
	}

	return &JobManager{jobDB: jobDB}, nil
}

// InsertData: Model의 InsertJobData 함수를 호출하여 JobData 삽입
func (jm *JobManager) InsertData(jobData *Model.JobData) error {
	err := jm.jobDB.InsertJobData(jobData)
	if err != nil {
		return err
	}
	return nil
}

func (jm *JobManager) popData(agentUUID string) (*Model.JobData, error, bool) {
	job, err, exist := jm.getDataByAgentUUID(agentUUID)
	if err != nil {
		return nil, err, exist
	}
	if exist == false {
		return nil, fmt.Errorf("Error! NoRecord"), exist
	}

	return job, nil, exist

}

// GetDataByAgentUUID: Model의 GetJobDataByAgentUUID 함수를 호출하여 JobData 조회
func (jm *JobManager) getDataByAgentUUID(agentUUID string) (*Model.JobData, error, bool) {
	job, err, exist := jm.jobDB.SelectJobDataByAgentUUID(agentUUID)
	if err != nil {
		return nil, err, exist
	}

	return job, nil, exist
}

// DeleteDataByInstructionUUID: Model의 DeleteJobDataByInstructionUUID 함수를 호출하여 JobData 삭제
func (jm *JobManager) deleteDataByInstructionUUID(instructionUUID string) error {
	err := jm.jobDB.DeleteJobDataByInstructionUUID(instructionUUID)
	if err != nil {
		return err
	}
	return nil
}
